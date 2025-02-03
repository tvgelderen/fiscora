import { PUBLIC_API_URL } from "$env/static/public";
import { redirect, type RequestHandler } from "@sveltejs/kit";

export const GET: RequestHandler = async ({ locals: { session } }) => {
	if (session) {
		return redirect(302, "/profile");
	}

	return redirect(302, `${PUBLIC_API_URL}/auth/google`);
};
