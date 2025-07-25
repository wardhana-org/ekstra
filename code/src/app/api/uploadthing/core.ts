import { createUploadthing, type FileRouter } from "uploadthing/next";
import { UploadThingError } from "uploadthing/server";
import { auth } from "~/server/auth";

const f = createUploadthing();

// FileRouter for your app, can contain multiple FileRoutes
export const ourFileRouter = {
  // Define as many FileRoutes as you like, each with a unique routeSlug
  imageUploader: f({
    image: {
      /**
       * For full list of options and defaults, see the File Route API reference
       * @see https://docs.uploadthing.com/file-routes#route-config
       */
      maxFileSize: "4MB",
      maxFileCount: 1,
    },
  })
    // Set permissions and file types for this FileRoute
    .middleware(async () => {
      // This code runs on your server before upload
      const user = await auth();

      // If you throw, the user will not be able to upload
      if (!user) throw new UploadThingError("Unauthorized");

      // Whatever is returned here is accessible in onUploadComplete as `metadata`
      return { userId: user.user.id };
    })
    .onUploadComplete(async ({ metadata, file }) => {
      // This code RUNS ON YOUR SERVER after upload
      console.log("Upload complete for userId:", metadata.userId);

      console.log("file url", file.ufsUrl);

      // !!! Whatever is returned here is sent to the clientside `onClientUploadComplete` callback
      return { uploadedBy: metadata.userId };
    }),

  postContentuploader: f({
    image: {maxFileSize: "4MB", maxFileCount: 1},
    video: {maxFileSize: "16MB", maxFileCount: 1},
    audio: {maxFileSize :"8MB", maxFileCount: 1},
    blob: {maxFileSize: "8MB", maxFileCount: 1},
    pdf: {maxFileSize: "4MB", maxFileCount: 1},
    text: {maxFileSize: "64KB", maxFileCount: 1}

  })
    .middleware(async () => {
      const user = await auth();
      const creatorPageId = user?.user.creatorPageId ?? null

      if (!user) throw new UploadThingError("Unauthorized");
      if (creatorPageId === null) throw new UploadThingError("Unauthorized");
      return { userId: user.user.id, creatorPageId: creatorPageId };

    })
    .onUploadComplete(async ({ metadata, file }) => {
      console.log("Upload complete for userId:", metadata.userId);
      console.log("file url", file.ufsUrl);

      return 
    }),

} satisfies FileRouter;

export type OurFileRouter = typeof ourFileRouter;
