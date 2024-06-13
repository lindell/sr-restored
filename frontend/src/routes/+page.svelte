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

<h1>SR restored</h1>

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
		Detta projekt har som mål att hjälpa Sveriges Radio att uppfylla deras uppdrag genom att
		generera Podcast RSS flöden med allt tillgängligt innehåll.
	</section>

	<section>
		<h2>Ett litet manifest... typ</h2>

		Podcasts är en underbar öppen teknologi där det öppna gränssnittet (podcast RSS flöden) tillåter
		både producenterna av media, samt lyssningsverktygen, att vara oberoende av varandra ❤️ Bland de
		stängda plattformarna som vi sett växa fram de senaste åren, så står podcasts fortfarande kvar
		som en av de få riktigt öppna teknologierna. Dessvärre har flera streamingplatformar börjat låsa
		innehåll exklusivt till deras tjänster. Detta är tråkigt, men med privata bolag så är det
		förståeligt, och något som är lätt att bojkotta om man inte tycker om metoderna. Att public
		service artificiellt börjar låsa innehåll är dock förbryllande. Snälla Sveriges Radio, gör så
		detta projektet kan arkiveras genom att inte exklusivt börja lansera innehåll på Sveriges Radio
		Play.

		<section>
			Snälla Sveriges Radio, gör så detta projektet kan arkiveras genom att inte exlusivt börja
			lansera innehåll på Sveriges Radio Play.
		</section>
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
		font-size: min(4rem, 15vw);
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
