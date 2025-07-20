## [SR Restored](https://sr-restored.se)

Sveriges Radio började 2023 plocka bort innehåll från sina RSS-flöden och andra plattformar, för att exklusivt lansera innehåll på SR Play.

[Detta går helt emot SRs definierade uppdrag](https://sverigesradio.se/artikel/vart-uppdrag)

> Sveriges Radios uppdrag är att leverera oberoende journalistik och kulturupplevelser till publiken där den finns och kan lyssna.

**Detta projekt har som mål att hjälpa Sveriges Radio att uppfylla sitt uppdrag genom att generera Podcast-RSS-flöden med allt tillgängligt innehåll.**

### Ett litet manifest... typ

Podcasts är en underbar öppen teknologi där det öppna gränssnittet (podcast-RSS-flöden) tillåter både producenterna av media, samt lyssningsverktygen, att vara oberoende av varandra ❤️ Bland de stängda plattformarna som vi sett växa fram de senaste åren, så står podcasts fortfarande kvar som en av de få riktigt öppna teknologierna. Dessvärre har flera streamingplattformar börjat låsa innehåll exklusivt till sina egna tjänster. Detta är tråkigt, men med privata bolag så är det förståeligt, och något som är lätt att bojkotta om man inte tycker om metoderna. Att public service artificiellt börjar låsa innehåll är dock förbryllande. Snälla Sveriges Radio, gör så att det här projektet kan arkiveras genom att inte exklusivt börja lansera innehåll på Sveriges Radio Play.

Snälla Sveriges Radio, gör så detta projektet kan arkiveras genom att inte exklusivt börja lansera innehåll på Sveriges Radio Play.

# Tekniskt om projektet

Projektets huvudsakliga fokus är att vara lätt och billigt att hosta. Cachning körs därför direkt i minnet istället för t.ex. Redis. Den kör för tillfället på en maskin med 256 MB RAM + 1 vCPU och klarar att av tiotals miljoner requests per dag med den konfigurationen.

# Hur använder jag projektet?

1. Hitta URLen till podcasten: [sr-restored.se](https://sr-restored.se)
2. Skriv in URLen i valfri podcastspelare
3. Nu kan du fortsätta lyssna på podden i din favoritspelare 🎉
