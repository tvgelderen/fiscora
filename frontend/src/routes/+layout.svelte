<script>
	import "../app.css";
	import { onNavigate } from "$app/navigation";
	import { Toaster } from "$lib/components/ui/sonner";
	import Header from "$lib/components/header.svelte";
	import Footer from "$lib/components/footer.svelte";

	let { children } = $props();

	onNavigate((navigation) => {
		if (!document.startViewTransition) return;

		return new Promise((resolve) => {
			document.startViewTransition(async () => {
				resolve();
				await navigation.complete;
			});
		});
	});
</script>

<Toaster position="bottom-center" />

<Header />

<main class="mx-auto min-h-[calc(100vh_-_var(--header-height))] w-full max-w-screen-xl px-2 pb-16 md:px-4">
	<div
		class="absolute inset-0 z-[-1] [background:radial-gradient(125%_125%_at_50%_15%,#00000000_40%,#3c14ffbb_200%)]"
	></div>
	{@render children()}
</main>

<Footer />
