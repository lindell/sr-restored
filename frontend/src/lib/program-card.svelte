<script lang="ts">
	import { PUBLIC_BASE_URL } from '$env/static/public';

	import type { Program } from './types/program';

	export let program: Program;
	$: rssUrl = `${PUBLIC_BASE_URL}/rss/${program.id}`;

	function copy() {
		navigator.clipboard.writeText(rssUrl);
	}

	function select(e: MouseEvent) {
		if (e.target instanceof Element) {
			window.getSelection()?.selectAllChildren(e.target);
		}
	}
</script>

<div class="program">
	<div class="thumbnail">
		<img src={program.image} alt={program.name} class="thumbnail-image" width="512" height="512" />
	</div>

	<div class="content">
		<h2>{program.name}</h2>
		<div class="description">
			{program.description}
		</div>
	</div>

	<div class="link">
		<code on:click={select}>{rssUrl}</code>

		<div>
			<button class="copy-button" on:click={copy}>Kopiera</button>
		</div>
	</div>
</div>

<style lang="scss">
	@use 'sass:color';

	$border-radius: 2rem;
	$thumbnail-offset: 2rem;

	.program {
		position: relative;
		display: inline-block;
		border-radius: $border-radius;
		margin: 1rem;
		margin-top: $thumbnail-offset + 1rem;
		background: #eef4ed;
		color: #252323;
		box-shadow: 1px 1px 5px #00000011;
		flex: 0 0 25rem;
		max-width: 100%;
	}

	.thumbnail {
		width: 94%;
		margin: 0 auto;
		border-radius: $border-radius;
		margin-top: -$thumbnail-offset;
		overflow: hidden;
		box-shadow: 1px 1px 5px #00000063;
	}

	.thumbnail img {
		width: 100%;
		height: auto;
		display: block;
	}

	h2 {
		text-align: center;
		font-weight: 400;
		text-decoration: underline;
		padding: 0.5em 0;
		margin: 0;
	}

	.description {
		text-align: center;
	}

	.content {
		text-align: center;
		padding: 0.25rem 1rem;
	}

	.link {
		text-align: center;
		width: 100%;
	}

	.link code {
		display: block;
	}

	.link > * {
		margin: 1rem;
	}

	.copy-button {
		background-color: #18a972;
		color: #fff;
		padding: 1em 4em;
		border: none;
		border-radius: 99rem;
		cursor: pointer;
	}

	.copy-button:hover {
		background-color: color.adjust(#18a972, $lightness: -5%);
	}

	.copy-button:active {
		background-color: color.adjust(#18a972, $lightness: -10%);
	}
</style>
