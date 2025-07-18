import "server-only";
import { db } from "~/server/db";
import { eq } from "drizzle-orm";
import * as schema from "~/server/db/schema";

// session data query
// find if user is a creator
export async function getCreatorId(userId: string) {
  const creatorPage = await db.query.creatorPages.findFirst({
    where: eq(schema.creatorPages.userId, userId),
  });
  return creatorPage?.id ?? null;
}

// check handle availability
export async function isHandleAvailable(handle: string) {
  const existing = await db.query.creatorPages.findFirst({
    where: eq(schema.creatorPages.pageHandle, handle),
  });
  return !existing; // Returns true if handle is available
}

// create creator page form queries
export async function createCreatorPage(
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
