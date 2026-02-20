<script lang="ts">
	import infoCircle from './info-circle.svg';

	export let label: string;
	export let info: string;

	const id = `info-popover-${Math.random().toString(36).slice(2, 9)}`;
	let open = false;

	function toggle() {
		open = !open;
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape' && open) {
			open = false;
			e.stopPropagation();
		}
	}

	function handleBlur(e: FocusEvent) {
		const wrapper = (e.currentTarget as HTMLElement);
		// Delay to allow focus to move within the wrapper
		setTimeout(() => {
			if (!wrapper.contains(document.activeElement)) {
				open = false;
			}
		}, 0);
	}
</script>

<!-- svelte-ignore a11y_no_static_element_interactions -->
<div class="info-wrapper" on:keydown={handleKeydown} on:focusout={handleBlur}>
	<button
		class="info-btn"
		aria-label="Info om {label}"
		aria-expanded={open}
		aria-controls={id}
		on:click={toggle}
	>
		<img src={infoCircle} alt="" class="info-icon" />
	</button>
	<div
		{id}
		class="feed-popover"
		class:open
		role="tooltip"
		aria-hidden={!open}
	>
		<strong>{label}</strong>
		<p>{info}</p>
	</div>
</div>

<style lang="scss">
	.info-wrapper {
		position: relative;
	}

	.info-btn {
		background: none;
		border: none;
		cursor: pointer;
		padding: 0;
		line-height: 1;
		color: #555;

		&:hover {
			color: #222;
		}
	}

	.info-icon {
		width: 1.2rem;
		height: 1.2rem;
		display: block;
		opacity: 0.5;

		&:hover {
			opacity: 0.8;
		}
	}

	.feed-popover {
		position: absolute;
		top: calc(100% + 8px);
		left: 50%;
		transform: translateX(-50%);
		padding: 1rem 1.25rem;
		border-radius: 0.75rem;
		border: 1px solid #ccc;
		box-shadow: 0 4px 16px #0002;
		width: max-content;
		max-width: 20rem;
		font-size: 0.9rem;
		background: #fff;
		z-index: 10;
		display: none;

		strong {
			display: block;
			margin-bottom: 0.25rem;
		}

		p {
			margin: 0;
		}

		&::after {
			content: '';
			position: absolute;
			top: -8px;
			left: 50%;
			margin-left: -7px;
			width: 14px;
			height: 14px;
			background: #fff;
			border-left: 1px solid #ccc;
			border-top: 1px solid #ccc;
			transform: rotate(45deg);
		}
	}

	.feed-popover.open {
		display: block;
	}
</style>
