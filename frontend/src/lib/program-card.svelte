<script lang="ts">
	import { PUBLIC_BASE_URL } from '$env/static/public';

	import type { Program } from './types/program';

	export let program: Program;
	$: rssUrl = `${PUBLIC_BASE_URL}/rss/${program.id}`;

	function copy() {
		navigator.clipboard.writeText(rssUrl);
	}
</script>

<div class="program">
	<div class="thumbnail">
		<img src={program.image} alt={program.name} class="thumbnail" width="512" height="512" />
	</div>

	<div class="content">
		<h2>{program.name}</h2>
		<div class="description">
			{program.description}
		</div>

		<div class="link">
			<code>{rssUrl}</code>

			<div>
				<button class="copy-button" on:click={copy}>Kopiera</button>
			</div>
		</div>
	</div>
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
		color: #000;
	}

	.thumbnail img {
		width: 100%;
		height: auto;
	}

	h2 {
		text-align: center;
	}

	.description {
		text-align: center;
	}

	.content {
		text-align: center;
		padding: 1rem;
	}

	.link {
		margin-top: 1.5rem;
	}
</style>
