<script lang="ts">
	import JobStatus from '$lib/jobStatus';
	import { onMount } from 'svelte';
	import JobComponent from './JobComponent.svelte';

	interface Notification {
		body: string;
		status: JobStatus;
		sent: boolean;
		name: string;
	}

	let messages: Notification[] = [];
	// let socket;

	const addMessage = (message: string) => {
		messages = [...messages, JSON.parse(message)];
	};

	onMount(() => {
		const socket = new WebSocket('ws://localhost:8090/ws');

		socket.addEventListener('open', function (event) {
			console.log('websocket connection opened.');
		});

		socket.addEventListener('message', function (event) {
			addMessage(event.data);
			console.log('received message:', event.data);
		});

		return () => {
			console.log('websocket connection closed.');
			socket.close();
		};
	});
</script>

<div class="flex flex-col my-4 gap-2">
	{#each messages as message}
		<JobComponent jobTitle={message.name} jobStatus={message.status} jobBody={message.body} />
	{/each}
</div>

<style lang="postcss">
	:global(html) {
		background-color: theme(colors.gray.100);
	}
</style>
