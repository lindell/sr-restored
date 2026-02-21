import { XMLParser } from 'fast-xml-parser';

interface XMLProgram {
	'@_id': string;
	'@_name': string;
	description: string;
	programimage: string;
	programurl: string;
}

let cached: ReturnType<typeof doFetch> | null = null;

async function doFetch() {
	const xml = await fetch('https://api.sr.se/api/v2/programs?size=10000');

	const parser = new XMLParser({
		ignoreAttributes: false,
		attributeNamePrefix: '@_'
	});
	const parsedXML = parser.parse(await xml.text());
	const xmlPrograms: XMLProgram[] = parsedXML.sr.programs.program;

	return xmlPrograms.map((program) => ({
		id: Number(program['@_id']),
		name: program['@_name'],
		description: program.description,
		image: program.programimage,
		url: program.programurl
	}));
}

export function fetchPrograms() {
	if (!cached) {
		cached = doFetch();
	}
	return cached;
}
