<script lang="ts">
	import ProgramCard from '$lib/program-card/program-card.svelte';
	import type { Program } from '$lib/types/program';

	export let data: PageData;

	const program: Program = data.program;
</script>

<h1><a href="/">SR restored</a></h1>

<div class="program-detail">
	<ProgramCard {program} full includeFeeds />
</div>

<style>
	h1 {
		margin-top: min(5vw, 5vh);
		text-align: center;
		font-size: min(4rem, 15vw);

		a:visited {
			color: inherit;
		}

		a {
			text-decoration: none;
			color: inherit;
		}
	}

	.program-detail {
		display: flex;
		justify-content: center;
		margin-top: 2rem;
	}
</style>
