import { error } from '@sveltejs/kit';

export const prerender = false;
export const ssr = false;

export async function load({ params, parent }) {
	const parentData = await parent();

	const program = parentData.programs.find((p) => p.id === Number(params.id));

	if (!program) error(404);

	return {
		program
	};
}
