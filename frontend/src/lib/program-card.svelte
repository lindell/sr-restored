<script lang="ts">
	import { PUBLIC_BASE_URL } from '$env/static/public';

	import type { Program } from './types/program';

	export let program: Program;
</script>

<div class="program">
	<div class="thumbnail">
		<img src={program.image} alt={program.name} class="thumbnail" width="512" height="512" />
	</div>
	<h2>{program.name}</h2>
	<div class="description">
		{program.description}
	</div>

	{`${PUBLIC_BASE_URL}/rss/${program.id}`}
</div>

<style lang="scss">
	.program {
		max-width: 30rem;
		width: 100%;
		display: inline-block;
		border-radius: 1em;
		overflow: hidden;
		margin: 1em;
		background: #ddb967;
	}

	.thumbnail img {
		width: 100%;
	}

	h2 {
		text-align: center;
	}

	.description {
		text-align: center;
	}
</style>
