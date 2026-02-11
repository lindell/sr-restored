<script lang="ts">
	import { PUBLIC_BASE_URL } from '$env/static/public';
	import Button from './button.svelte';

	import type { Program } from './types/program';

	export let program: Program;
	export let full: boolean = false;
	export let includeFeeds: boolean = false;

	function copy(text: string) {
		navigator.clipboard.writeText(text);
	}

	function select(e: MouseEvent) {
		if (e.target instanceof Element) {
			window.getSelection()?.selectAllChildren(e.target);
		}
	}
</script>

<div
	class="program {full ? 'full' : ''}"
	style="view-transition-name: program-{program.id}; view-transition-class: program-card;"
>
	<div class="thumbnail" style="view-transition-name: program-thumb-{program.id};">
		<a href={'/programs/' + program.id}
			><img
				src={program.image}
				alt={program.name}
				class="thumbnail-image"
				width="512"
				height="512"
			/></a
		>
	</div>

	<div class="content">
		<h2>
			<a style="view-transition-name: program-title-{program.id};" href={'/programs/' + program.id}
				>{program.name}</a
			>
		</h2>
		<div class="description">
			{program.description}
		</div>
	</div>

	{#if includeFeeds}
		<div class="link">
			<!-- svelte-ignore a11y_click_events_have_key_events a11y_no_noninteractive_element_interactions -->
			{#each ['', '/broadcast', '/download'] as suffix}
				{@const url = `${PUBLIC_BASE_URL}/rss/${program.id}${suffix}`}
				<div class="feed-row">
					<div>Normala RSS flödet, fungerar för de flesta.</div>
					<div class="feed-link">
						<code on:click={select}>{url}</code>
						<div class="copy-btn">
							<Button on:click={() => copy(url)}>Kopiera</Button>
						</div>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>

<style lang="scss">
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
		max-width: 50rem;
	}

	.program.full {
		flex: 0 0 100%;
	}

	.program:hover .thumbnail {
		transform: scale(1.02);
	}

	.thumbnail {
		width: 94%;
		margin: 0 auto;
		border-radius: $border-radius;
		margin-top: -$thumbnail-offset;
		overflow: hidden;
		box-shadow: 1px 1px 5px #00000063;
		max-width: 40vh;
		transition: transform 0.2s;
	}

	.thumbnail img {
		width: 100%;
		height: auto;
		display: block;
	}

	h2 {
		text-align: center;
		font-weight: 400;
		padding: 0.5em 0;
		margin: 0;

		a {
			color: inherit;
		}
	}

	.description {
		text-align: center;
	}

	.content {
		text-align: center;
		padding: 1.5rem 1rem;
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

	.feed-link {
		display: flex;
		justify-content: center;
		align-items: center;
		gap: 1rem;
		flex-wrap: wrap;
	}

	.copy-btn {
		margin-top: 0.5rem;
	}
</style>
