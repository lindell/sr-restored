import { error } from '@sveltejs/kit';
import { fetchPrograms } from '$lib/programs';

export const prerender = true;

export async function entries() {
	const programs = await fetchPrograms();
	return programs.map((p) => ({ id: String(p.id) }));
}

export async function load({ params }) {
	const programs = await fetchPrograms();
	const program = programs.find((p) => p.id === Number(params.id));

	if (!program) error(404);

	return {
		program
	};
}
