import { redirect } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";
import { getTransactionsYearInfo, getTransactionsYearInfoPerType } from "$lib/api/transactions";
import { getYear } from "$lib/api/utils";

export const load: PageServerLoad = async ({ locals: { session, user }, url }) => {
	if (!session?.accessToken || !user) {
		throw redirect(302, "/login");
	}

	const year = getYear(url.searchParams);

	return {
		yearInfo: await getTransactionsYearInfo(year, session.accessToken),
		incomeInfo: await getTransactionsYearInfoPerType(year, true, session.accessToken),
		expenseInfo: await getTransactionsYearInfoPerType(year, false, session.accessToken),
	};
};
