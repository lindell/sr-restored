<script lang="ts">
	import ProgramCard from '$lib/program-card.svelte';
	import type { Program } from './types/program';

	export let programs: Program[];

	let searchQuery = '';

	$: results = programs.filter((program) =>
		program.name.toLowerCase().includes(searchQuery.toLowerCase())
	);

	function onSearchFocus(e: FocusEvent) {
		if (e.target instanceof Element) {
			e.target.scrollIntoView({ behavior: 'smooth', block: 'start' });
		}
	}
</script>

<div class="search">
	<input
		class="search-input"
		type="text"
		placeholder="Sök program"
		bind:value={searchQuery}
		on:focus={onSearchFocus}
	/>
</div>

<div class="programs">
	{#each results.slice(0, 12) as program (program.id)}
		<ProgramCard {program} />
	{/each}
	{#if results.length === 0}
		Inga resultat matchar "{searchQuery}"
	{/if}
</div>

<style lang="scss">
	.programs {
		display: flex;
		flex-wrap: wrap;
		justify-content: center;
		align-items: center;
	}

	.search {
		display: flex;
		justify-content: center;
	}

	.search-input {
		border: none;
		outline: solid 1.5px #9e9e9e;
		font-size: 1.5rem;
		border-radius: 1.5rem;
		background: none;
		padding: 1.5rem;
		color: #fff;
		width: 30rem;
		max-width: calc(100% - 3rem);
		scroll-margin: 1em;
	}

	.search-input:focus {
		outline: solid 2px #e3e3e3;
	}
</style>
