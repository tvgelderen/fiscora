import {
	getExpenseTypes,
	getIncomeTypes,
	getTransactionIntervals,
	getTransactions,
	getTransactionsYearInfo,
} from "$lib/api/transactions";
import { redirect } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";
import { getMonth, getYear } from "$lib/api/utils";

export const load: PageServerLoad = async ({ locals: { session, user }, url }) => {
	if (!session?.accessToken || !user) {
		throw redirect(302, "/login");
	}

	const year = getYear(url.searchParams);
	const month = getMonth(url.searchParams);

	return {
		transactions: await getTransactions(month, year, session.accessToken),
		transactionIntervals: await getTransactionIntervals(session.accessToken),
		incomeTypes: await getIncomeTypes(session.accessToken),
		expenseTypes: await getExpenseTypes(session.accessToken),
		yearInfo: await getTransactionsYearInfo(year, session.accessToken),
		demo: user.isDemo,
	};
};
