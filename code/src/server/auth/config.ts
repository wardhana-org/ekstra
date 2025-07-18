import { DrizzleAdapter } from "@auth/drizzle-adapter";
import type { DefaultSession, NextAuthConfig } from "next-auth";
import GoogleProvider from "next-auth/providers/google";
import { getCreatorId } from "~/server/queries";
import { db } from "~/server/db";
import {
	accounts,
	sessions,
	users,
	verificationTokens,
} from "~/server/db/schema";
import { eq } from "drizzle-orm";


/**
 * Module augmentation for `next-auth` types. Allows us to add custom properties to the `session`
 * object and keep type safety.
 *
 * @see https://next-auth.js.org/getting-started/typescript#module-augmentation
 */
declare module "next-auth" {
	interface Session extends DefaultSession {
		user: {
			id: string;
      role: string;
      creatorPageId: string | null;
			// ...other properties
			// role: UserRole;
		} & DefaultSession["user"];
	}

  interface User {
    role: string;
  }
}

/**
 * Options for NextAuth.js used to configure adapters, providers, callbacks, etc.
 *
 * @see https://next-auth.js.org/configuration/options
 */
export const authConfig = {
	providers: [
		GoogleProvider({
    clientId: process.env.AUTH_GOOGLE_ID,
    clientSecret: process.env.AUTH_GOOGLE_SECRET
  })
		/**
		 * ...add more providers here.
		 *
		 * Most other providers require a bit more work than the Discord provider. For example, the
		 * GitHub provider requires you to add the `refresh_token_expires_in` field to the Account
		 * model. Refer to the NextAuth.js docs for the provider you want to use. Example:
		 *
		 * @see https://next-auth.js.org/providers/github
		 */
	],
	adapter: DrizzleAdapter(db, {
		usersTable: users,
		accountsTable: accounts,
		sessionsTable: sessions,
		verificationTokensTable: verificationTokens,
	}),
	callbacks: {
    // one session per user check
    async signIn({ user }) {
      if (user?.id) {
        await db.delete(sessions).where(eq(sessions.userId, user.id));
      }
      return true;
    },
		session: async ({ session, user }) => ({
      ...session,
      user: {
        ...session.user,
        id: user.id,
        role: user.role,
        creatorPageId: await getCreatorId(user.id)
      },
		}),
	},
} satisfies NextAuthConfig;
