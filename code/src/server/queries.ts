import "server-only";
import { db } from "~/server/db";
import { eq, sql } from "drizzle-orm";
import * as schema from "~/server/db/schema";

// session data query
// find if user is a creator
export async function getCreatorId(userId: string) {
  const creatorPage = await db.query.creatorPages.findFirst({
    where: eq(schema.creatorPages.userId, userId),
  });
  return creatorPage?.id ?? null;
}

// creator queries
// check handle availability
export async function isHandleAvailable(handle: string) {
  const existing = await db.query.creatorPages.findFirst({
    where: eq(schema.creatorPages.pageHandle, handle),
  });
  return !existing; // Returns true if handle is available
}

// create creator page form queries
export async function createMyCreatorPage(
  data: {
    name: string;
    handle: string;
    description: string;
    userId: string;
    userImage: string | null;
  }
) {
  const newCreatorPage = await db.insert(schema.creatorPages).values({
    name: data.name,
    userId: data.userId,
    description: data.description,
    pageHandle: data.handle,
    profileImage: data.userImage,
  });
  return newCreatorPage;
}

// fetch creator account data (for account dropdown so profile image and page name) based on creator page id from session
export async function getMyCreatorPageData(
  data: {
    creatorPageId: string;
  }
) {
  const creatorPageData = 
    await db.select({
      name: schema.creatorPages.name,
      description: schema.creatorPages.description,
      profileImage: schema.creatorPages.profileImage,
      pageHandle: schema.creatorPages.pageHandle,
    })
    .from(schema.creatorPages)
    .where(eq(schema.creatorPages.id, data.creatorPageId))
    .limit(1);

  return creatorPageData
}

export async function getMyCreatorDropdownData(
  data: {
    creatorPageId: string;
  }
) {
  const creatorDropdownData = 
    await db.select({
      name: schema.creatorPages.name,
      profileImage: schema.creatorPages.profileImage,
    })
    .from(schema.creatorPages)
    .where(eq(schema.creatorPages.id, data.creatorPageId))
    .limit(1);

  return creatorDropdownData;
}

export async function createMyMembershipTier( data: {
  title: string;
  price: number;
  description: string;
  creatorPageId: string;

}) {
  const newMembershipTier = 
    await db.insert(schema.memberships).values({
      creatorPageId: data.creatorPageId,
      title: data.title,
      description: data.description,
      price: sql`${data.price}::numeric(10,2)`
    });
  return newMembershipTier
}

export async function getMyMembershipTiers(data: {
  creatorPageId : string
} ) {
  const myMembershipTiers =
    await db.select({
      id: schema.memberships.id,
      title: schema.memberships.title,
      description: schema.memberships.description,
      price: schema.memberships.price,
    })
    .from(schema.memberships)
    .where(eq(schema.memberships.creatorPageId, data.creatorPageId))
  return myMembershipTiers
}

export async function createMyContent(data: {
  creatorPageId: string,
  type: string,
  contentKey: string,
  usedIn: string,
}) {
  const myContent =
    await db.insert(schema.contents).values({
      creatorPageId: data.creatorPageId,
      type: data.type,
      contentKey: data.contentKey,
      usedIn: data.usedIn
    }).returning({ id: schema.contents.id });
  return myContent
}

export async function createMyPost(data: {
  title: string,
  description: string,
  creatorPageId: string,
}) {
  const myPost = 
    await db.insert(schema.posts).values({
      title: data.title,
      description: data.description,
      creatorPageId: data.creatorPageId,
    }).returning({ id: schema.posts.id });
  return myPost
}

export async function createMyPostContent (
  postId: string,
  contentId: string
) {
  const myPostContent =
    await db.insert(schema.postContents).values({
      contentId: contentId,
      postId: postId
    }).returning({ id: schema.postContents.id });
  return myPostContent
}

export async function createMyMembershipExclusivePost( data: {
  postId: string,
  membershipId: string,
}) {
  const membershipExclusivePost =
    await db.insert(schema.membershipContents).values({
      postId: data.postId,
      membershipId: data.membershipId,
    }).returning({ id: schema.membershipContents.id });
  return membershipExclusivePost
}

export async function getMyMembershipName( data : {
  creatorPageId : string
}) {
  const myMembershipNames =
    await db.select({
      id: schema.memberships.id,
      title: schema.memberships.title,
    })
    .from(schema.memberships)
    .where(eq(schema.memberships.creatorPageId, data.creatorPageId))
  return myMembershipNames
}

// app queries
// images queries
export async function getBackground() {
  const backgroundImages = await db.query.appImages.findMany({
    columns: {
      url: true,
    },
    orderBy: (model, { desc }) => [desc(model.id)],
    limit: 3,
  });
  return backgroundImages;
}
