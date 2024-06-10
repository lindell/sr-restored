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

<input type="text" placeholder="Sök" bind:value={searchQuery} />

<div class="programs">
	{#each results.slice(0, 12) as program (program.id)}
		<ProgramCard {program} />
	{/each}
</div>

<style lang="scss">
	h1 {
		color: #fff;
	}

	.programs {
		display: flex;
		flex-wrap: wrap;
		justify-content: space-around;
	}
</style>
