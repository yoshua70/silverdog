<script lang="ts">
	let url = '';

	function handleSubmit() {
		// Handle the form submission here
		console.log('Submitted URL:', url);
	}
</script>

<form on:submit|preventDefault={handleSubmit} class="max-w-md w-full">
	<div class="mb-4">
		<label for="url" class="block text-gray-700 mb-2">Enter an URL to download from:</label>
		<input
			type="url"
			id="url"
			bind:value={url}
			class="shadow appearance-none border rounded w-full py-3 px-4 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
			placeholder="Enter URL"
			required
		/>
	</div>

	<div class="flex items-center justify-between">
		<button
			type="submit"
			class="bg-pink-300 hover:bg-pink-400 text-gray-900 font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline w-full"
		>
			Submit
		</button>
	</div>
</form>
