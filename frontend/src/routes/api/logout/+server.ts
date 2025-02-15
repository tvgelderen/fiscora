import { PUBLIC_API_URL } from "$env/static/public";
import { forbidden } from "$lib";
import { redirect, type RequestHandler } from "@sveltejs/kit";

export const GET: RequestHandler = async ({ locals: { session } }) => {
	if (!session) {
		return forbidden();
	}

	return redirect(302, `${PUBLIC_API_URL}/auth/logout`);
};
