<script lang="ts">
	import Plus from "lucide-svelte/icons/plus";
	import TransactionList from "./transaction-list.svelte";
	import TransactionInfo from "./transaction-info.svelte";
	import TransactionForm from "./transaction-form.svelte";
	import TransactionMonthHeader from "./transaction-month-header.svelte";
	import { IncomingTypes, type Transaction, type TransactionMonthInfo } from "../../../ambient";
	import { getCurrentMonthNumber, getCurrentYear, listAllMonths } from "$lib";
	import * as Dialog from "$lib/components/ui/dialog";
	import * as Popover from "$lib/components/ui/popover/index.js";
	import { Button, buttonVariants } from "$lib/components/ui/button";
	import { CalendarIcon } from "lucide-svelte";
	import MonthPicker from "$lib/components/month-picker.svelte";
	import { onMount, tick } from "svelte";
	import { cn } from "$lib/utils";

	let { data } = $props();
	let { transactions, transactionIntervals, incomeTypes, expenseTypes, yearInfo, demo } = data;

	let showFormModal: boolean = $state(false);
	let year: number = $state(getCurrentYear());
	let month: number = $state(getCurrentMonthNumber());
	let incoming = $state(IncomingTypes[0]);
	let transactionsState: Transaction[] = $state(transactions);
	let monthInfo: TransactionMonthInfo | null = $state(null);
	let monthInfoDiff: TransactionMonthInfo | null = $state(null);
	let selectedTransaction: Transaction | null = $state(null);
	let editTransaction: Transaction | null = $state(null);

	const months = listAllMonths();

	function setSelectedTransaction(transaction: Transaction | null) {
		selectedTransaction = transaction;
	}

	async function setEditTransaction(transaction: Transaction | null) {
		editTransaction = transaction;
		await tick();
		showFormModal = true;
	}

	function closeFormModal() {
		editTransaction = null;
		showFormModal = false;
	}

	async function updateTransactions() {
		const response = await fetch(`/api/transactions?month=${month}&year=${year}`);
		const value = (await response.json()) as Transaction[];
		transactionsState = value;
	}

	function updateMonthInfo() {
		const info = yearInfo.get(month.toString());
		if (!info) return;

		monthInfo = info;

		const prevMonth = yearInfo.get((month - 1).toString());
		if (!prevMonth) return;

		monthInfoDiff = {
			income: monthInfo.income - prevMonth.income,
			expense: monthInfo.expense - prevMonth.expense,
		};
	}

	async function handleMonthChanged() {
		await updateTransactions();
		updateMonthInfo();
	}

	async function handleSuccess() {
		closeFormModal();
		await updateTransactions();
	}

	function addTransaction(transaction: Transaction, idx: number) {
		transactionsState.splice(idx, 0, transaction);
		transactionsState = [...transactionsState];
		updateYearInfo(transaction.amount, true);
	}

	function removeTransaction(transaction: Transaction) {
		transactionsState = transactionsState.filter((t) => t.id !== transaction.id);
		updateYearInfo(transaction.amount, false);
	}

	function updateYearInfo(amount: number, add: boolean) {
		if (monthInfo === null) return;

		if (add) {
			if (amount > 0) {
				monthInfo.income += amount;
			} else {
				monthInfo.expense -= amount;
			}
		} else {
			if (amount > 0) {
				monthInfo.income -= amount;
			} else {
				monthInfo.expense += amount;
			}
		}

		yearInfo.set(month.toString(), monthInfo);
		updateMonthInfo();
	}

	$effect(() => {
		transactionsState = transactions.filter((t) => {
			return (
				incoming === IncomingTypes[0] ||
				(incoming === IncomingTypes[1] && t.amount > 0) ||
				(incoming === IncomingTypes[2] && t.amount < 0)
			);
		});
	});

	onMount(() => {
		updateMonthInfo();
	});
</script>

<svelte:head>
	<title>Fiscora - Transactions</title>
</svelte:head>

<div class="mx-auto mb-8 text-center lg:mb-12">
	<h1 class="mb-4">Your transactions</h1>
	<p>Add, view, and edit your transactions to stay on top of your financial journey.</p>
	<p>Track your finances with ease and gain valuable insights.</p>
</div>

<TransactionMonthHeader {monthInfo} {monthInfoDiff} />

<div class="flex flex-wrap items-center justify-center gap-4 sm:flex-nowrap sm:justify-between sm:gap-6">
	<div class="order-last flex gap-2 sm:order-first">
		{#each IncomingTypes as incomingType}
			<button
				class="rounded-full px-4 py-2 backdrop-blur-[1px] transition-colors {incoming !== incomingType &&
					'hover:bg-primary/10'} {incoming === incomingType && 'bg-primary/10 ring-[1px] ring-primary/50'}"
				onclick={() => (incoming = incomingType)}
			>
				{incomingType}
			</button>
		{/each}
	</div>

	<Popover.Root>
		<Popover.Trigger
			class={cn("!h-fit w-[240px] justify-start text-left text-base", buttonVariants({ variant: "outline" }))}
		>
			<CalendarIcon class="mr-2 h-5 w-5" />
			{months.get(month)}
		</Popover.Trigger>
		<Popover.Content class="w-auto p-0">
			<MonthPicker bind:year bind:month callback={handleMonthChanged} />
		</Popover.Content>
	</Popover.Root>
</div>

<TransactionList
	transactions={transactionsState}
	select={setSelectedTransaction}
	edit={setEditTransaction}
	add={addTransaction}
	remove={removeTransaction}
	{demo}
/>

<Dialog.Root
	bind:open={showFormModal}
	onOpenChange={(open) => {
		if (!open) {
			editTransaction = null;
		}
	}}
>
	<TransactionForm
		{transactionIntervals}
		{incomeTypes}
		{expenseTypes}
		{demo}
		transaction={editTransaction}
		success={handleSuccess}
		close={closeFormModal}
	/>
</Dialog.Root>

<Dialog.Root
	open={selectedTransaction !== null}
	onOpenChange={(open) => {
		if (!open) {
			selectedTransaction = null;
		}
	}}
>
	<TransactionInfo transaction={selectedTransaction} />
</Dialog.Root>

<button
	class="fixed bottom-5 right-4 rounded-full bg-primary p-2 text-slate-50 shadow-md !shadow-slate-500 transition-all duration-300 hover:shadow-lg dark:!shadow-slate-800 sm:bottom-8 sm:right-8"
	onclick={() => (showFormModal = true)}
	disabled={demo}
>
	<Plus size={24} />
	<span class="sr-only">Add Budget</span>
</button>
