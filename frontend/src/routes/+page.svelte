<script lang="ts">
	import Search from '$lib/search.svelte';
	import type { PageData } from './$types';

	export let data: PageData;
</script>

<h1>SR restored</h1>

<div class="info">
	<section>
		Sveriges Radio började 2023 plocka bort innehåll från sina RSS-flöden och andra plattformar, för
		att exklusivt lansera innehåll på SR Play.
	</section>

	<section>
		<a href="https://sverigesradio.se/artikel/vart-uppdrag"
			>Detta går helt emot Sveriges Radios definierade uppdrag!</a
		>
	</section>

	<section class="qoute">
		“Sveriges Radios uppdrag är att leverera oberoende journalistik och kulturupplevelser till
		publiken där den finns och kan lyssna.”
	</section>

	<section>
		Detta projekt har som mål att hjälpa Sveriges Radio att uppfylla sitt uppdrag genom att
		generera Podcast-RSS-flöden med allt tillgängligt innehåll.
	</section>

	<section>
		<h2>Ett litet manifest... typ</h2>

		Podcasts är en underbar öppen teknologi där det öppna gränssnittet (podcast-RSS-flöden) tillåter
		både producenterna av media, samt lyssningsverktygen, att vara oberoende av varandra ❤️ Bland de
		stängda plattformarna som vi sett växa fram de senaste åren, så står podcasts fortfarande kvar
		som en av de få riktigt öppna teknologierna. Dessvärre har flera streamingplattformar börjat låsa
		innehåll exklusivt till sina egna tjänster. Detta är tråkigt, men med privata bolag så är det
		förståeligt, och något som är lätt att bojkotta om man inte tycker om metoderna. Att public
		service artificiellt börjar låsa innehåll är dock förbryllande. Snälla Sveriges Radio, gör så
		att det här projektet kan arkiveras genom att inte exklusivt börja lansera innehåll på Sveriges Radio
		Play.

		<section>
			Snälla Sveriges Radio, gör så att det här projektet kan arkiveras genom att inte exklusivt börja
			lansera innehåll på Sveriges Radio Play.
		</section>
	</section>

	<section class="signing">
		<a
			href="https://www.skrivunder.com/sveriges_radio_bor_sluta_publicera_innehall_exlusivt_pa_egna_plattformar#form"
			>Skriv under för att Sveriges Radio bör börja publicera allt innehåll för alla
			podcastanvändare igen!</a
		>
	</section>
</div>

<hr />

<Search programs={data.programs} />

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

	.signing a {
		color: #18a972;
		font-weight: bold;
	}
</style>
