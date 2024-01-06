<script>
	import '../app.css';
	import Sidebar from '../components/Sidebar.svelte';
</script>

<div class="flex">
	<Sidebar />
	<main
		class="flex flex-col items-center gap-8 bg-gray-100 rounded-xl p-4 my-12 mx-4 w-full min-h-[700px]"
	>
		<h1 class="text-4xl text-center font-semibold">silverdog</h1>

		<slot />
	</main>
</div>
