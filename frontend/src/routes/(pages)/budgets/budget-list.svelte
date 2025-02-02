<script lang="ts">
	import type { Budget } from "../../../ambient";
	import { Edit, Trash } from "lucide-svelte";
	import * as AlertDialog from "$lib/components/ui/alert-dialog";
	import { getFormattedAmount, getFormattedDateShortWithYear } from "$lib";
	import { toast } from "svelte-sonner";
	import { buttonVariants } from "$lib/components/ui/button";

	let {
		budgets,
		demo,
		add,
		edit,
		remove,
	}: {
		budgets: Budget[];
		demo: boolean;
		add: (idx: number, budget: Budget) => void;
		edit: (budget: Budget) => void;
		remove: (budget: Budget) => void;
	} = $props();

	let budgetToDelete: Budget | null = $state(null);

	function openDeleteConfirmation(budget: Budget) {
		budgetToDelete = budget;
	}

	function closeDeleteModal() {
		budgetToDelete = null;
	}

	async function deleteBudget() {
		if (budgetToDelete === null) return;

		const id = budgetToDelete.id;
		const idx = budgets.findIndex((b) => b.id === id);
		const budget = budgets.at(idx);

		remove(budgetToDelete);

		closeDeleteModal();

		if (demo) {
			toast.warning("Demo users cannot delete budgets");
			return;
		}

		const response = await fetch(`/api/budgets/${id}`, { method: "DELETE" });
		if (!response.ok) {
			toast.error("Error deleting budget");

			if (budget !== undefined) {
				add(idx, budget);
			}

			return;
		}

		toast.success("Budget deleted successfully");
	}
</script>

<div class="flex flex-wrap justify-center gap-6">
	{#each budgets as budget (budget.id)}
		<div class="card flex w-full max-w-sm flex-col justify-between p-6">
			<div>
				<a href="/budgets/{budget.id}" class="hover:underline">
					<h3 class="mb-2">
						{budget.name}
					</h3>
				</a>
				<p class="mb-4 text-sm text-muted-foreground">
					{budget.description}
				</p>
				<div class="mb-4 flex flex-col gap-2 rounded-lg bg-muted-foreground/10 p-3">
					<div class="flex items-center justify-between">
						<span class="text-sm font-medium">Total Budget:</span>
						<span class="text-lg font-bold">
							{getFormattedAmount(budget.amount)}
						</span>
					</div>
					<div class="flex justify-between text-sm text-muted-foreground">
						<span>
							Start: {getFormattedDateShortWithYear(budget.startDate)}
						</span>
						<span>
							End: {getFormattedDateShortWithYear(budget.endDate)}
						</span>
					</div>
				</div>
				<div class="mb-2">
					<ul class="list-inside list-disc space-y-2">
						{#each budget.expenses as expense}
							<li class="flex items-center justify-between">
								<span class="text-base text-muted-foreground">
									{expense.name}:
								</span>
								<span class="font-semibold">
									{getFormattedAmount(expense.allocatedAmount)}
								</span>
							</li>
						{/each}
					</ul>
				</div>
			</div>
			<div class="mt-4 flex justify-end gap-1">
				<button
					class={buttonVariants({ size: "icon", variant: "ghost" })}
					onclick={() => edit(budget)}
					disabled={demo}
				>
					<Edit size={20} />
				</button>
				<button
					class={`${buttonVariants({ size: "icon", variant: "ghost" })} hover:bg-destructive`}
					onclick={() => openDeleteConfirmation(budget)}
					disabled={demo}
				>
					<Trash size={20} />
				</button>
			</div>
		</div>
	{/each}
</div>

{#if budgets.length === 0}
	<p class="text-center text-muted-foreground">You haven't set any budgets yet. Create one to get started!</p>
{/if}

<AlertDialog.Root open={budgetToDelete !== null}>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<h3>Confirm Deletion</h3>
		</AlertDialog.Header>
		<p class="mb-4">Are you sure you want to delete this budget? This action is permanent and cannot be undone.</p>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={closeDeleteModal}>Cancel</AlertDialog.Cancel>
			<AlertDialog.Action class={buttonVariants({ variant: "destructive" })} onclick={deleteBudget}>
				Delete
			</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
