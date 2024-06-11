<script lang="ts">
	import ProgramCard from '$lib/program-card.svelte';
	import type { PageData } from './$types';

	let searchQuery = '';

	export let data: PageData;

	$: results = data.programs.filter((program) =>
		program.name.toLowerCase().includes(searchQuery.toLowerCase())
	);

	function onSearchFocus(e: FocusEvent) {
		if (e.target instanceof Element) {
			e.target.scrollIntoView({ behavior: 'smooth', block: 'start' });
		}
	}
</script>

<h1>Sverige Radio Unsensored</h1>

<div class="info">
	<section>
		Sveriges Radio började 2023 plocka bort innehåll från sina RSS flöden och andra platformar, för
		att exlusivt lansera innehåll på SR Play.
	</section>

	<section>
		<a href="https://sverigesradio.se/artikel/vart-uppdrag"
			>Detta går helt emot SRs definierade uppdrag!</a
		>
	</section>

	<section class="qoute">
		“Sveriges Radios uppdrag är att leverera oberoende journalistik och kulturupplevelser till
		publiken där den finns och kan lyssna.”
	</section>

	<section>
		Denna tjänst tillhandahåller Podcast RSS flöden som inkluderar det bortblockade innehållet. Så
		att du kan använda den podcast-spelaren du föredrar.
	</section>
</div>

<hr />

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
</div>

<style lang="scss">
	h1 {
		margin-top: min(10vw, 15vh);
		text-align: center;
	}

	hr {
		border: none;
		background-color: #414141;
		height: 0.125rem;
		margin: 2rem;
	}

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

	.info {
		max-width: 50em;
		text-align: center;
		margin-left: auto;
		margin-right: auto;
		margin-bottom: 2rem;
	}

	.info a {
		color: inherit;
	}

	.info section {
		margin: 1rem 0;
	}

	.qoute {
		border-left: solid 0.125rem #b3b3b3;
		padding: 0 0.5em;
	}
</style>
