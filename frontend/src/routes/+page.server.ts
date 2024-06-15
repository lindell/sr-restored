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
	// Just popular
	[5067, 96], // P3 Historia
	[909, 95], // P1 Dokumentär
	[4845, 94], // Creepypodden i P3
	[516, 93], // Spanarna
	[5413, 92], // P3 Krim
	[4540, 91], // Ekot nyhetssändning
	[5419, 90], // Fråga Agnes Wold
	[4941, 89] // Europapodden
]);

interface XMLProgram {
	'@_id': string;
	'@_name': string;
	description: string;
	programimage: string;
	programurl: string;
}
