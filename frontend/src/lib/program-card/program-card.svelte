<script lang="ts">
	import { PUBLIC_BASE_URL } from '$env/static/public';
	import Button from '../button.svelte';
	import InfoPopover from './info-popover.svelte';
	import caretUp from './caret-up.svg';
	import caretDown from './caret-down.svg';

	import type { Program } from '../types/program';

	export let program: Program;
	export let full: boolean = false;
	export let includeFeeds: boolean = false;

	let showMoreFeeds = false;
	let link = full ? program.url : '/programs/' + program.id;

	const feeds = [
		{
			suffix: '',
			label: 'Standard',
			info: 'Det vanliga RSS-flödet. Funkar för de flesta. Använd om du är osäker på vilket flöde du behöver.'
		},
		{
			suffix: '/broadcast',
			label: 'Broadcast',
			info: 'RSS-flöde med ljud från livesändningar. Använd om du vill lyssna på programmet som det sänts, inte som det publicerats i efterhand.'
		},
		{
			suffix: '/on-demand',
			label: 'On demand',
			info: 'RSS-flöde med ljud från klippta avsnitt. Använd om du vill lyssna på programmet som det publicerats i efterhand, inte som det sänts.'
		}
	];

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
		<a href={link}
			><img
				src={program.image}
				alt={program.name}
				class="thumbnail-image"
				width="512"
				height="512"
				fetchpriority="high"
			/></a
		>
	</div>

	<div class="program-inner">
		<div class="content">
			<h2>
				<a style="view-transition-name: program-title-{program.id};" href={link}>{program.name}</a>
			</h2>
			<div class="description">
				{program.description}
			</div>
		</div>

		{#if includeFeeds}
			<div class="link">
				<!-- svelte-ignore a11y_click_events_have_key_events a11y_no_noninteractive_element_interactions -->
				<div class="feed-grid">
					{#each feeds as feed, i (feed.label)}
						{@const url = `${PUBLIC_BASE_URL}/rss/${program.id}${feed.suffix}`}
						{@const hidden = i > 0 && !showMoreFeeds}
						<div class="feed-row-item" class:hidden>
							<InfoPopover label={feed.label} info={feed.info} />
						</div>
						<code class:hidden on:click={select}>{url}</code>
						<div class="feed-row-item copy-btn" class:hidden>
							<Button on:click={() => copy(url)}>Kopiera</Button>
						</div>
					{/each}
				</div>

				<button class="toggle-feeds" on:click={() => (showMoreFeeds = !showMoreFeeds)}>
					{showMoreFeeds ? 'Dölj alternativa flöden' : 'Visa alternativa flöden'}
					<img src={showMoreFeeds ? caretUp : caretDown} alt="" class="caret-icon" />
				</button>
			</div>
		{/if}
	</div>
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

	.program-inner {
		margin: 1rem;
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
		padding: 0 1rem;
	}

	.link > * {
		margin: 1rem;
	}

	.feed-grid {
		display: inline-grid;
		grid-template-columns: auto auto auto;
		column-gap: 1rem;
		align-items: center;
		margin-top: 0;
	}

	.feed-grid > :nth-child(n + 4) {
		margin-top: 0.75rem;
	}

	.feed-grid code {
		text-align: center;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
		direction: rtl;
	}

	.hidden {
		visibility: hidden;
		height: 0;
		overflow: hidden;
		margin: 0 !important;
		padding: 0;
		border: none;
	}

	.toggle-feeds {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 0.4rem;
		margin-left: auto;
		margin-right: auto;
		background: none;
		border: none;
		color: #555;
		cursor: pointer;
		font-size: 0.9rem;
		padding: 0.5rem 1rem;
		margin-top: 0;
		margin-bottom: 0.5rem;

		&:hover {
			color: #222;
		}
	}

	.caret-icon {
		width: 0.75rem;
		height: 0.75rem;
		opacity: 0.6;
	}
</style>
