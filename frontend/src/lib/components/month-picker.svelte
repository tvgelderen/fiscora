<script lang="ts">
	import { ChevronLeft, ChevronRight } from "lucide-svelte";
	import { getCurrentMonthNumber, getCurrentYear, listAllMonthsShort } from "$lib";

	let {
		year = $bindable(),
		month = $bindable(),
		callback,
	}: {
		year: number;
		month: number;
		callback?: () => void;
	} = $props();

	let yearState = $state(year);

	const months = listAllMonthsShort();
	const currentYear = getCurrentYear();
	const currentMonth = getCurrentMonthNumber();

	function previousYear() {
		yearState--;
	}

	function nextYear() {
		yearState++;
	}

	function select(newYear: number, newMonth: number) {
		year = newYear;
		month = newMonth;
		if (callback) {
			callback();
		}
	}
</script>

<div class="w-[240px] p-4">
	<div class="flex flex-col space-y-4">
		<div class="flex w-full items-center justify-between">
			<button
				name="previous-year"
				aria-label="Go to previous year"
				class="inline-flex h-7 w-7 items-center justify-center whitespace-nowrap rounded-md border border-input bg-transparent p-0 text-sm font-medium opacity-50 ring-offset-background transition-colors hover:bg-accent hover:text-accent-foreground hover:opacity-100 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50"
				onclick={previousYear}
			>
				<ChevronLeft class="h-4 w-4" />
			</button>
			<div class="text-sm font-medium" aria-live="polite" role="presentation" id="month-picker">
				{yearState}
			</div>
			<button
				name="next-year"
				aria-label="Go to next year"
				class="inline-flex h-7 w-7 items-center justify-center whitespace-nowrap rounded-md border border-input bg-transparent p-0 text-sm font-medium opacity-50 ring-offset-background transition-colors hover:bg-accent hover:text-accent-foreground hover:opacity-100 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50"
				onclick={nextYear}
			>
				<ChevronRight class="h-4 w-4" />
			</button>
		</div>
		<div class="grid w-full grid-cols-3 gap-2" role="grid" aria-labelledby="month-picker">
			{#each months.entries() as [k, v]}
				<div
					class="relative p-0 text-center text-sm focus-within:relative focus-within:z-20 [&:has([aria-selected])]:bg-slate-100 first:[&:has([aria-selected])]:rounded-l-md last:[&:has([aria-selected])]:rounded-r-md dark:[&:has([aria-selected])]:bg-slate-800"
					role="presentation"
				>
					<button
						class={`inline-flex h-9 w-16 items-center justify-center rounded-md p-0 text-sm transition-colors hover:bg-slate-100 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-slate-400 focus-visible:ring-offset-2 dark:ring-offset-slate-950 dark:hover:bg-slate-800 dark:hover:text-slate-50 dark:focus-visible:ring-slate-800 ${k === month && year === yearState && "bg-primary text-slate-50 hover:bg-primary focus:bg-primary dark:hover:bg-primary dark:focus:bg-primary"} ${k === currentMonth && k !== month && yearState === currentYear && "bg-slate-100 dark:bg-slate-800"}`}
						role="gridcell"
						tabIndex={-1}
						onclick={() => select(yearState, k)}
					>
						{v}
					</button>
				</div>
			{/each}
		</div>
	</div>
</div>
