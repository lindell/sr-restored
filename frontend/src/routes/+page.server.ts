export const prerender = true;

import { XMLParser } from 'fast-xml-parser';
import type { PageLoad } from './$types';

export const load: PageLoad = async () => {
	const xml = await fetch('https://api.sr.se/api/v2/programs?size=10000');

	const parser = new XMLParser({
		ignoreAttributes: false,
		attributeNamePrefix: '@_'
	});
	const parsedXML = parser.parse(await xml.text());
	const xmlPrograms: XMLProgram[] = parsedXML.sr.programs.program;

	const programs = xmlPrograms
		.map((program) => ({
			id: Number(program['@_id']),
			name: program['@_name'],
			description: program.description,
			image: program.programimage,
			url: program.programurl
		}))
		.sort((a, b) => (ranks.get(b.id) ?? 0) - (ranks.get(a.id) ?? 0));

	return {
		programs
	};
};

// Used to reorder frontpage
const ranks = new Map<number, number>([
	// Known to limit content:
	[2519, 100], // P3 Dokumentär
	[2071, 99], // Sommar & Vinter i P1
	[2024, 98], // Morgonpasset i P3
	[4923, 97], // USApodden
	[4845, 96], // Creepypodden i P3
	[5419, 95], // Fråga Agnes Wold
	// Just popular
	[5067, 94], // P3 Historia
	[909, 93], // P1 Dokumentär
	[516, 92], // Spanarna
	[5413, 91], // P3 Krim
	[4540, 90], // Ekot nyhetssändning
	[4941, 89], // Europapodden
	[5188, 88], // P3 Dystopia
	[4947, 87], // P3 Serie
	[5524, 86] // Ekonomiakuten
]);

interface XMLProgram {
	'@_id': string;
	'@_name': string;
	description: string;
	programimage: string;
	programurl: string;
}
