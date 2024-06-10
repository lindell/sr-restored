<script lang="ts">
	import ProgramCard from '$lib/program-card.svelte';
	import type { PageData } from './$types';

	let searchQuery = '';

	export let data: PageData;

	$: results = data.programs.filter((program) =>
		program.name.toLowerCase().includes(searchQuery.toLowerCase())
	);
</script>

<h1>Sverige Radio Unsensored</h1>

<div class="search">
	<input class="search-input" type="text" placeholder="Sök program" bind:value={searchQuery} />
</div>

<hr />

<div class="programs">
	{#each results.slice(0, 12) as program (program.id)}
		<ProgramCard {program} />
	{/each}
</div>

<style lang="scss">
	h1 {
		margin-top: 5em;
		text-align: center;
		color: #fff;
	}

	hr {
		border-color: #9e9e9e;
		margin: 2rem;
	}

	.programs {
		display: flex;
		flex-wrap: wrap;
		justify-content: space-around;
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
	}

	.search-input:focus {
		outline: solid 2px #e3e3e3;
	}
</style>
