<script lang="ts">
	import { onMount } from "svelte";
	import X from "lucide-svelte/icons/x";
	import Sun from "lucide-svelte/icons/sun";
	import Moon from "lucide-svelte/icons/moon";
	import Menu from "lucide-svelte/icons/menu";
	import User from "lucide-svelte/icons/user";
	import { createDarkMode } from "$lib/theme.svelte";

	type NavLink = {
		link: string;
		title: string;
	};

	let navLinks: NavLink[] = [
		{
			link: "/",
			title: "Home",
		},
		{
			link: "/transactions",
			title: "Transactions",
		},
		{
			link: "/dashboard",
			title: "Dashboard",
		},
		{
			link: "/budgets",
			title: "Budgets",
		},
		{
			link: "/reports",
			title: "Reports",
		},
	];

	let navOpen = $state(false);

	const darkMode = createDarkMode();

	const toggleNav = () => (navOpen = !navOpen);
	const closeNav = () => (navOpen = false);

	const toggleTheme = () => {
		darkMode.toggle();
		const theme = darkMode.darkMode ? "dark" : "light";
		const html = document.querySelector("html");
		localStorage.setItem("theme", theme);
		if (html) {
			html.classList.value = theme;
		}
	};

	onMount(() => {
		window.addEventListener("resize", () => {
			if (window.innerWidth >= 1024) {
				closeNav();
			}
		});

		const theme = localStorage.getItem("theme");
		if (theme) {
			darkMode.set(theme === "dark");
		} else {
			const prefersDark = window.matchMedia("(prefers-color-scheme: dark)");
			if (prefersDark) {
				darkMode.set(true);
			} else {
				darkMode.set(false);
			}
		}
	});
</script>

<header class="relative z-10 flex h-[var(--header-height)] w-full items-center justify-between px-4">
	<a href="/" class="hover:text-current">
		<h1 class="text-4xl">Fiscora.</h1>
	</a>

	<nav class="absolute left-[50%] flex h-full translate-x-[-50%] items-center">
		<ul class="text-md hidden h-full items-center gap-6 lg:flex">
			{#each navLinks as link}
				<li>
					<a class="hover-underline" href={link.link} aria-label={link.title}>
						{link.title}
					</a>
				</li>
			{/each}
		</ul>
	</nav>

	<div class="flex items-center">
		<button onclick={toggleTheme} class="icon mr-4" id="theme-toggle">
			{#if darkMode.darkMode}
				<Sun />
			{:else}
				<Moon />
			{/if}
		</button>
		<a href="/profile"><User class="hidden lg:block" /></a>
		<button class="block lg:hidden" onclick={toggleNav} aria-label="menu">
			<Menu size={32} />
		</button>
	</div>
</header>

<!-- Side navbar -->
{#if navOpen}
	<div class="absolute inset-0 z-[100] bg-background/25 backdrop-blur-sm" onclick={closeNav} role="none"></div>
{/if}
<div
	class="absolute bottom-0 {navOpen
		? 'left-0'
		: 'left-[-400px]'} top-0 z-[100] w-full max-w-[400px] bg-background transition-all duration-300"
>
	<button class="absolute right-2 top-2" onclick={closeNav} aria-label="close-nav">
		<X />
	</button>
	<ul id="side-nav" class="flex h-dvh w-full flex-col items-center justify-center gap-4 text-xl">
		{#each navLinks as link}
			<li>
				<a href={link.link} onclick={closeNav}>
					{link.title}
				</a>
			</li>
		{/each}
	</ul>
</div>
<!-- Side navbar -->
