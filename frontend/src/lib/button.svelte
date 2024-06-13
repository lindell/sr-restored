<button on:click {...$$restProps}><slot /></button>

<style lang="scss">
	@use 'sass:color';

	button {
		background-color: #18a972;
		color: #fff;
		padding: 1em 4em;
		border: none;
		border-radius: 99rem;
		cursor: pointer;
	}

	button:hover {
		background-color: color.adjust(#18a972, $lightness: -5%);
	}

	button:active {
		background-color: color.adjust(#18a972, $lightness: -10%);
	}
</style>
