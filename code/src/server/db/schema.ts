import { relations, sql } from "drizzle-orm";
import { index, pgTableCreator, primaryKey } from "drizzle-orm/pg-core";
import type { AdapterAccount } from "next-auth/adapters";

/**
 * This is an example of how to use the multi-project schema feature of Drizzle ORM. Use the same
 * database instance for multiple projects.
 *
 * @see https://orm.drizzle.team/docs/goodies#multi-project-schema
 */
export const createTable = pgTableCreator((name) => `ekstra_${name}`);

// App stuff
export const appImages = createTable("appImage", 
  (d) => ({
    id: d.integer().primaryKey().generatedByDefaultAsIdentity(),
    name: d.varchar({ length: 256 }).notNull(),
    url: d.varchar({ length: 1024 }).notNull(),
    createdAt: d
      .timestamp({ withTimezone: true })
      .default(sql`CURRENT_TIMESTAMP`)
      .notNull(),
    updatedAt: d.timestamp({ withTimezone: true }).$onUpdate(() => new Date()),
  }),
  (t) => [

	],
)

// user functionality
export const users = createTable("user",
  (d) => ({
    id: d
      .varchar({ length: 255 })
      .notNull()
      .primaryKey()
      .$defaultFn(() => crypto.randomUUID()),
    name: d.varchar({ length: 255 }),
    email: d.varchar({ length: 255 }).notNull(),
    emailVerified: d
      .timestamp({
        mode: "date",
        withTimezone: true,
      })
      .default(sql`CURRENT_TIMESTAMP`),
    role:d
      .varchar({length:50})
      .notNull()
      .default("user")
      .$type<"user" | "admin">(),
    image: d.varchar({ length: 255 }),
}));

export const creatorPages = createTable(
  "creatorPage", 
  (d) => ({
    id: d
      .varchar({ length: 255 })
      .notNull()
      .primaryKey()
      .$defaultFn(() => crypto.randomUUID()),
    name: d
      .varchar({ length: 255 }).notNull(),
    userId: d
      .varchar({ length: 255 })
      .notNull()
      .references(() => users.id),
    description: d
      .varchar({ length: 255 }),
    pageHandle: d
      .varchar({ length: 1024 }).notNull().unique(), // user created
    profileImage: d
      .varchar({ length: 255 }), // default image should be taken from user profile image
    createdAt: d
      .timestamp({ withTimezone: true, mode: 'date' })
      .notNull()
      .default(sql`CURRENT_TIMESTAMP`),
    updatedAt: d
      .timestamp({ withTimezone: true, mode: 'date' })
  }),
  (t) =>[
    index("creator_page_user_id_idx").on(t.userId),
    index("creator_page_id_idx").on(t.id),
  ]
);

export const contents = createTable(
  "content",
  (d) => ({
    id: d
      .varchar({ length: 255 })
      .notNull()
      .primaryKey()
      .$defaultFn(() => crypto.randomUUID()),
    creatorPageId: d
      .varchar({ length: 255 })
      .notNull()
      .references(() => creatorPages.id),
    type: d // video, image, file, etc.
      .varchar({ length: 100})
      .notNull(),
    contentKey: d // upload thing
      .varchar({ length: 1024 }).notNull(),
    usedIn: d
      .varchar({ length: 25 }),
  	createdAt: d
			.timestamp({ withTimezone: true })
			.default(sql`CURRENT_TIMESTAMP`)
			.notNull(),
		updatedAt: d.timestamp({ withTimezone: true }).$onUpdate(() => new Date()),
  }),
  (t) => [

  ]
)

export const storeListings = createTable (
  "storeListing",
  (d) => ({
    id: d
      .varchar({ length: 255 })
      .notNull()
      .primaryKey()
      .$defaultFn(() => crypto.randomUUID()),
    title: d
      .varchar({ length: 50 }).notNull(),
    description: d
      .varchar({ length: 50 }).notNull(),
    price: d
      .numeric({ precision: 10, scale: 2 })
      .notNull(),
    creatorPageId: d
      .varchar({ length: 255 })
      .notNull()
      .references(() => creatorPages.id),
  	createdAt: d
			.timestamp({ withTimezone: true })
			.default(sql`CURRENT_TIMESTAMP`)
			.notNull(),
		updatedAt: d.timestamp({ withTimezone: true }).$onUpdate(() => new Date()),
  }),
  (t) => [

  ]
)

export const storeContents = createTable(
  "storeContent",
  (d) => ({
    id: d
      .varchar({ length: 255 })
      .notNull()
      .primaryKey()
      .$defaultFn(() => crypto.randomUUID()),
    contentId: d
			.varchar({ length: 255 })
			.notNull()
			.references(() => contents.id),
    storeListingId: d
      .varchar({ length: 255 })
			.notNull()
			.references(() => storeListings.id),
  })
)

export const posts = createTable(
  "post",
  (d) => ({
    id: d
      .varchar({ length: 255 })
      .notNull()
      .primaryKey()
      .$defaultFn(() => crypto.randomUUID()),
    title: d
      .varchar({ length: 255 }).notNull(),
    description: d
      .varchar({ length: 255 }),
    creatorPageId: d
      .varchar({ length: 255 })
      .notNull()
      .references(() => creatorPages.id),
    createdAt: d
			.timestamp({ withTimezone: true })
			.default(sql`CURRENT_TIMESTAMP`)
			.notNull(),
		updatedAt: d.timestamp({ withTimezone: true }).$onUpdate(() => new Date()),
  }),
  (t) => [

  ]
)

export const postContents = createTable(
  "postContent",
  (d) => ({
    id: d
      .varchar({ length: 255 })
      .notNull()
      .primaryKey()
      .$defaultFn(() => crypto.randomUUID()),
		contentId: d
			.varchar({ length: 255 })
			.notNull()
			.references(() => contents.id),
    postId: d
			.varchar({ length: 255 })
			.notNull()
			.references(() => posts.id),
  }),
  (t) => [

  ]
)

export const memberships = createTable(
  "membership",
  (d) => ({
    id: d
      .varchar({ length: 255 })
      .notNull()
      .primaryKey()
      .$defaultFn(() => crypto.randomUUID()),
    creatorPageId: d
      .varchar({ length: 255 })
      .notNull()
      .references(() => creatorPages.id),
    title: d
      .varchar({ length: 255 }).notNull(),
    description: d
      .varchar({ length: 255 }),
    price: d
      .numeric({ precision: 10, scale: 2 })
      .notNull(),
    createdAt: d
			.timestamp({ withTimezone: true })
			.default(sql`CURRENT_TIMESTAMP`)
			.notNull(),
		updatedAt: d.timestamp({ withTimezone: true }).$onUpdate(() => new Date()),   
  }),
  (t) => [

  ]
)

export const membershipContents = createTable(
	"membershipContent",
	(d) => ({
		id: d
      .varchar({ length: 255 })
      .notNull()
      .primaryKey()
      .$defaultFn(() => crypto.randomUUID()),
		postId: d
			.varchar({ length: 255 })
			.notNull()
			.references(() => posts.id),
    membershipId: d
      .varchar({ length: 255 })
			.notNull()
			.references(() => memberships.id),
	}),
	(t) => [
    
	],
);

// auth
export const accounts = createTable(
	"account",
	(d) => ({
		userId: d
			.varchar({ length: 255 })
			.notNull()
			.references(() => users.id),
		type: d.varchar({ length: 255 }).$type<AdapterAccount["type"]>().notNull(),
		provider: d.varchar({ length: 255 }).notNull(),
		providerAccountId: d.varchar({ length: 255 }).notNull(),
		refresh_token: d.text(),
		access_token: d.text(),
		expires_at: d.integer(),
		token_type: d.varchar({ length: 255 }),
		scope: d.varchar({ length: 255 }),
		id_token: d.text(),
		session_state: d.varchar({ length: 255 }),
	}),
	(t) => [
		primaryKey({ columns: [t.provider, t.providerAccountId] }),
		index("account_user_id_idx").on(t.userId),
	],
);

export const accountsRelations = relations(accounts, ({ one }) => ({
	user: one(users, { fields: [accounts.userId], references: [users.id] }),
}));

export const sessions = createTable(
	"session",
	(d) => ({
		sessionToken: d.varchar({ length: 255 }).notNull().primaryKey(),
		userId: d
			.varchar({ length: 255 })
			.notNull()
			.references(() => users.id),
		expires: d.timestamp({ mode: "date", withTimezone: true }).notNull(),
	}),
	(t) => [index("t_user_id_idx").on(t.userId)],
);

export const sessionsRelations = relations(sessions, ({ one }) => ({
	user: one(users, { fields: [sessions.userId], references: [users.id] }),
}));

export const verificationTokens = createTable(
	"verificationToken",
	(d) => ({
		identifier: d.varchar({ length: 255 }).notNull(),
		token: d.varchar({ length: 255 }).notNull(),
		expires: d.timestamp({ mode: "date", withTimezone: true }).notNull(),
	}),
	(t) => [primaryKey({ columns: [t.identifier, t.token] })],
);

// relations
export const usersRelations = relations(users, ({ many, one }) => ({
	accounts: many(accounts),
  session: one(sessions),
  creatorPage: one(creatorPages),
}));

export const creatorPagesRelations = relations(creatorPages, ({ many, one }) => ({
  users: one(users, {
    fields: [creatorPages.userId],
    references: [users.id],
  }),
  contents: many(contents),
  storeListings: many(storeListings),
  posts: many(posts),
  memberships: many(memberships),
}));

export const contentsRelations = relations(contents, ({ one, many }) => ({
  creator: one(creatorPages, {
    fields: [contents.creatorPageId],
    references: [creatorPages.id],
  }),
  storeContents: many(storeContents),
  postContents: many(postContents),
}));

export const storeListingsRelations = relations(storeListings, ({ one, many }) => ({
  creator: one(creatorPages, {
    fields: [storeListings.creatorPageId],
    references: [creatorPages.id],
  }),
  storeContents: many(storeContents),
}));

export const storeContentsRelations = relations(storeContents, ({ one }) => ({
  content: one(contents, {
    fields: [storeContents.contentId],
    references: [contents.id],
  }),
  storeListing: one(storeListings, {
    fields: [storeContents.storeListingId],
    references: [storeListings.id],
  }),
}));

export const postsRelations = relations(posts, ({ one, many }) => ({
  creator: one(creatorPages, {
    fields: [posts.creatorPageId],
    references: [creatorPages.id],
  }),
  postContents: many(postContents),
}));

export const postContentsRelations = relations(postContents, ({ one }) => ({
  post: one(posts, {
    fields: [postContents.postId],
    references: [posts.id],
  }),
  content: one(contents, {
    fields: [postContents.contentId],
    references: [contents.id],
  }),
}));

export const membershipsRelations = relations(memberships, ({ one, many }) => ({
  creator: one(creatorPages, {
    fields: [memberships.creatorPageId],
    references: [creatorPages.id],
  }),
  membershipContents: many(membershipContents),
}));

export const membershipContentsRelations = relations(membershipContents, ({ one }) => ({
  membership: one(memberships, {
    fields: [membershipContents.membershipId],
    references: [memberships.id],
  }),
  post: one(posts, {
    fields: [membershipContents.postId],
    references: [posts.id],
  }),
}));

