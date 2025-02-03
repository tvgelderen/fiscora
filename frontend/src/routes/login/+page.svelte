<script lang="ts">
	import { CircleUser } from "lucide-svelte";
</script>

<div class="mt-6">
	<div class="mx-auto w-full max-w-lg rounded-md border border-muted p-8 text-center shadow-xl">
		<h2 class="mb-4 text-3xl">Welcome to Fiscora</h2>
		<p class="mb-8">Take control of your finances today</p>

		<div class="mb-8">
			<h3 class="mb-4 text-2xl">Sign in</h3>
			<a class="btn" href="/api/login/google">
				<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 50 50" fill="currentColor" class="h-6 w-6">
					<path
						d="M 25.996094 48 C 13.3125 48 2.992188 37.683594 2.992188 25 C 2.992188 12.316406 13.3125 2 25.996094 2 C 31.742188 2 37.242188 4.128906 41.488281 7.996094 L 42.261719 8.703125 L 34.675781 16.289063 L 33.972656 15.6875 C 31.746094 13.78125 28.914063 12.730469 25.996094 12.730469 C 19.230469 12.730469 13.722656 18.234375 13.722656 25 C 13.722656 31.765625 19.230469 37.269531 25.996094 37.269531 C 30.875 37.269531 34.730469 34.777344 36.546875 30.53125 L 24.996094 30.53125 L 24.996094 20.175781 L 47.546875 20.207031 L 47.714844 21 C 48.890625 26.582031 47.949219 34.792969 43.183594 40.667969 C 39.238281 45.53125 33.457031 48 25.996094 48 Z"
					></path>
				</svg>
				Sign in with Google
			</a>
		</div>

		<div class="relative">
			<hr class="mb-8 mt-10 border-muted-foreground" />
			<span class="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 bg-background px-4 text-sm">
				or
			</span>
		</div>

		<div>
			<h3 class="mb-4 text-2xl">Try a demo</h3>
			<a class="secondary btn" href="/api/login/demo">
				<CircleUser />
				Login as demo user
			</a>
			<p class="mt-4 text-sm text-muted-foreground">Experience Fiscora without creating an account</p>
		</div>

		<div class="mt-8 text-xs text-muted-foreground">
			By logging in, you agree to our
			<a href="/terms-of-service" class="underline hover:text-indigo-950 dark:hover:text-indigo-200">
				Terms of Service
			</a>
			and
			<a href="/privacy-policy" class="underline hover:text-indigo-950 dark:hover:text-indigo-200">
				Privacy Policy
			</a>
			.
		</div>
	</div>
</div>
