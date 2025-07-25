'use server';

import { auth } from '~/server/auth';
import * as queries from '~/server/queries';
import { utapi } from './uploadthing';

export async function checkHandleUnique(handle: string) {
  try{
    return await queries.isHandleAvailable(handle); // returns true if taken
  } catch (error) {
    return null
  }
}

export async function submitCreatorPage(data: {
  name: string;
  handle: string;
  description: string;
}) {
  const user = await auth();
  const userId = user?.user?.id;
  const userImage = user?.user?.image ?? null;

  if (!userId) {
    throw new Error('User not authenticated');
  }

  try {
    await queries.createMyCreatorPage({
      ...data,
      userId,
      userImage,
    });
    return { success: true };
  } catch (error) {
    return { success: false, error: 'Failed to create creator page' };
  }
}

export async function retrieveCreatorDropdownData() {
  const user = await auth();
  const creatorPageId = user?.user.creatorPageId
  if (!creatorPageId) {
    return null;
  }
  try {
    const creatorDropdownData = await queries.getMyCreatorDropdownData({creatorPageId});
    return creatorDropdownData
  } catch (error) {
    return null;
  }
}

export async function retrieveCreatorPageData() {
  const user = await auth();
  const creatorPageId = user?.user.creatorPageId
  if (!creatorPageId) {
    return null;
  }
  try {
    const creatorDropdownData = await queries.getMyCreatorPageData({creatorPageId});
    return creatorDropdownData
  } catch (error) {
    return null;
  }
}

export async function submitMembershipTierInfo( data:{
  title: string;
  price: number;
  description: string;
}) {

  const user = await auth();
  const creatorPageId = user?.user.creatorPageId
  if (!creatorPageId) {
    return null;
  }

  try {
    await queries.createMyMembershipTier({
      ...data,
      creatorPageId
    })
    return { success: true };
  } catch(error) {
    return { success: false, error: 'Failed to create creator page' };
  }
}

export async function retrieveMyMembershipTiers() {
  const user = await auth();
  const creatorPageId = user?.user.creatorPageId
  if (!creatorPageId) {
    return null;
  }

  try {
    const myMembershipTiers =
      await queries.getMyMembershipTiers({
        creatorPageId
      })
    return myMembershipTiers
  } catch(error) {
    return null
  }
}

export async function retrieveMembershipNames() {
  const user = await auth();
  const creatorPageId = user?.user.creatorPageId
  if (!creatorPageId) {
    return null;
  }
  try {
    const myMembershipNames =
      await queries.getMyMembershipName({
        creatorPageId
      })
    return myMembershipNames
  } catch(error) {
    return null
  }
}

export async function removeFileFromUt(
  key: string
) {
  await utapi.deleteFiles(key);
}

export async function submitPostData(title: string, description: string) {
  const user = await auth()
  const creatorPageId = user?.user.creatorPageId
  if (!creatorPageId) {
    return null;
  }
  try {
    const result = await queries.createMyPost({
      title,
      description,
      creatorPageId
    })
    if(!result[0]) throw new Error("Failed to create post");
    const postId = result[0].id;
    return postId
  } catch (error) {
    return null
  }
}

export async function submitContentData(
  data :{
    key: string;
    type: string;
    usedIn: string;}[],
) {
  const user = await auth()
  const creatorPageId = user?.user.creatorPageId
  if (!creatorPageId) {
    return null;
  }
  try {
    const contentIds = await Promise.all(
      data.map(async (content) => {
        const result = await queries.createMyContent({
          creatorPageId: creatorPageId,
          type: content.type,
          contentKey: content.key,
          usedIn: content.usedIn
        })
        if (!result[0]) throw new Error("Failed to create content");
        return result[0].id;
      })
    ) 
    return contentIds 
  } catch (error) {
    return null
  }
}

export async function linkPostContent(
  postId: string,
  contentIds: string[]
) {
  try {
    const postContentIds = await Promise.all(
      contentIds.map(async (contentId) => {
        const result = await queries.createMyPostContent(postId, contentId)
        if (!result[0]) throw new Error("Failed to link post content");
        return result[0].id
      })
    )
    return postContentIds
  } catch (error) {
    return null
  }
}

export async function linkPostToMembership(data: {
  postId: string,
  membershipId: string,
}) {
  try {
    const membershipexclusiveId = 
      await queries.createMyMembershipExclusivePost(data)
    return membershipexclusiveId[0]?.id
  } catch (error) {
    return null
  }
}