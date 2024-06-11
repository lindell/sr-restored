<svelte:head>
	<title>Sveriges Radio uncensored</title>
	<meta
		name="description"
		content="Sveriges Radios podcasts med innehåll som bara släpps på SR-Play"
	/>
</svelte:head>

<a href="https://github.com/lindell/sr-uncensored" class="ribbon"
	><img
		decoding="async"
		width="149"
		height="149"
		src="https://github.blog/wp-content/uploads/2008/12/forkme_right_white_ffffff.png?resize=149%2C149"
		class="attachment-full size-full"
		alt="Fork me on GitHub"
		loading="lazy"
		data-recalc-dims="1"
	/></a
>

<slot />

<style>
	:global(body) {
		margin: 1rem;
		font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
		background: #252323;
		color: #eef4ed;
	}

	.ribbon {
		position: absolute;
		right: 0;
		top: 0;
	}

	@media (max-width: 600px) {
		.ribbon {
			display: none;
		}
	}
</style>
