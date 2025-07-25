'use client'

import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import { UploadButton } from '~/utils/uploadthing';
import { useState } from 'react';
import { linkPostContent, linkPostToMembership, removeFileFromUt, submitContentData, submitPostData } from '~/server/actions';
import { useRouter } from 'next/navigation';

interface Memberships {
  id: string;
  title: string;
}

type UploadedFile = {
  key: string;
  name: string;
  size: number;
  type: string;
};

const formSchema = z.object({
  title: z.string().min(1, "Title is required"),
  description: z.string().min(10, "Description must be at least 10 characters"),
  membershipIds: z.array(z.string()).min(1, "Please select at least one availability option"),
})

type FormValues = z.infer<typeof formSchema>;

export default function CreatePostForm({ memberships }: { memberships: Memberships[] }) {
  
  const labelStyle = 'block pb-2';
  const inputStyle = 'w-full px-4 py-2 rounded-lg bg-stone-800 border-2 border-gray-400 focus:outline-none focus:ring-1 focus:ring-white';
  const errorStyle = 'mt-1 text-sm text-red-500';
  
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    setValue,
    watch
  } = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    mode: 'onChange',
    defaultValues: {
      membershipIds: ['free'] // Default to "Free to All" selected
    }
  });

  // Handle "Free to All" checkbox change
  const handleFreeToAllChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.checked) {
      setValue('membershipIds', ['free']); // Select only "Free to All"
    } else {
      setValue('membershipIds', []); // Deselect all
    }
  };

  function handleRemoveFile(key: string): void {
    try{
      removeFileFromUt(key)
      setUploadedFiles(prev => prev.filter(f => f.key !== key));
    }
    catch (error) {
      console.error("Deletion error:", error);
    }
  }

  const router = useRouter();
  const [uploadedFiles, setUploadedFiles] = useState<UploadedFile[]>([]);
  const [isProcessing, setIsProcessing] = useState(false);
  const isButtonDisabled = isSubmitting || isProcessing;

  const onSubmit = async (data: FormValues) => {
    setIsProcessing(true);
    try {
      // create post return postId
      const postId = await submitPostData(data.title, data.description)
      if (!postId) {
        console.error("Post creation failed.");
        return;
      }

      // Check if there is content
      if (uploadedFiles.length) {
        // write content (key) to database and return content id
        const contentData = uploadedFiles.map(file => ({
          key: file.key,
          type: file.type,
          usedIn: 'post'
          // name and size are excluded
        }));
        const contentId = await submitContentData(contentData)
        if (!contentId) {
          console.error("Content creation failed.");
          return;
        }

        // link content with post return postContent id
        await linkPostContent(postId, contentId)
      }

      // if membership exclusive content, link post to memberships
      if (!data.membershipIds.includes('free')) {
        for (const membershipId of data.membershipIds) {
          await linkPostToMembership({postId, membershipId});
        }
      }

      // redirect after finished
      router.push('/creator/dashboard/home')

    } catch (error) {
      console.error('Error creating post:', error);
      alert('Failed to create post');
      setIsProcessing(false);
    }
  }

  return(
    <div className="w-full p-4 bg-neutral-800 rounded-lg shadow-sm">
      <h2 className="text-2xl font-bold mb-6">Post Details</h2>
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
        
        {/* Membership choice - Always show "Free to All" */}
        <div>
          <label className="block pb-2">
            Availability
          </label>
          
          {/* Free to All option (always shown) */}
          <div className="flex items-center mb-3">
            <input
              type="checkbox"
              id="membership-free"
              value="free"
              checked={watch("membershipIds")?.includes('free')}
              onChange={handleFreeToAllChange}
              className="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500"
            />
            <label
              htmlFor="membership-free"
              className="ml-2 text-sm font-medium text-gray-300 hover:text-white cursor-pointer"
            >
              Free to All
            </label>
          </div>

          {/* Other memberships (only shown if they exist) */}
          {memberships.length > 0 && (
            <>
              <label className="block pb-2 text-sm text-gray-400">
                Or restrict to specific memberships:
              </label>
              <div className="space-y-2 pl-4 border-l-2 border-gray-700">
                {memberships.map((membership) => (
                  <div key={membership.id} className="flex items-center">
                    <input
                      type="checkbox"
                      id={`membership-${membership.id}`}
                      value={membership.id}
                      {...register("membershipIds")}
                      disabled={watch("membershipIds")?.includes('free')}
                      className="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500 disabled:opacity-50"
                    />
                    <label
                      htmlFor={`membership-${membership.id}`}
                      className="ml-2 text-sm font-medium text-gray-300 hover:text-white cursor-pointer"
                    >
                      {membership.title}
                    </label>
                  </div>
                ))}
              </div>
            </>
          )}

          {errors.membershipIds && (
            <p className={errorStyle}>{errors.membershipIds.message}</p>
          )}

          <p className="text-sm text-gray-400 mt-1">
            {watch("membershipIds")?.includes('free') 
              ? "This post will be available to everyone"
              : "Select specific memberships for this post"}
          </p>
        </div>

        {/* Post Title */}
        <div>
          <label className={labelStyle}>
            Title
          </label>
          <input
            {...register("title")}
            className={inputStyle}
            placeholder="Video title"
          />
          {errors.title && (
            <p className={errorStyle}>{errors.title.message}</p>
          )}
        </div>

        {/* Content */}
        {uploadedFiles.map((file) => (
          <div key={file.key} className="p-3 bg-gray-700 rounded-lg mb-2">
            <div className="flex justify-between items-start">
              <div>
                <p className="font-medium text-gray-100">{file.name}</p>
                <p className="text-xs text-gray-400">
                  {file.size} • {file.type}
                </p>
              </div>
              <button
                type='button'
                onClick={() => handleRemoveFile(file.key)}
                className="text-red-400 hover:text-red-300 text-sm"
              >
                Remove
              </button>
            </div>
          </div>
        ))}
        <div>
          <label className="block pb-2">Content</label>
            <div className={isButtonDisabled ? 'pointer-events-none opacity-50' : ''}>
              <UploadButton 
                endpoint={'postContentuploader'}
                onClientUploadComplete={(res) => {
                  if (!res) return;
                  setUploadedFiles(res.map(file => ({
                    key: file.key,
                    name: file.name,
                    size: file.size,
                    type: file.type
                  })));
                  console.log("Upload complete:", res);
                }}
              />
            </div>
        </div>

        {/* Description */}
        <div>
          <label className="block pb-2">Description</label>
          <textarea
            {...register("description")}
            className="w-full px-4 py-2 rounded-lg bg-stone-800 border-2 border-gray-400"
            rows={5}
          />
          {errors.description && (
            <p className={errorStyle}>{errors.description.message}</p>
          )}
        </div>

        {/* Submit Button */}
        <button
          type="submit"
          disabled={isSubmitting}
          className={`mt-6 w-full py-2 px-4 text-black rounded-lg flex items-center justify-center transition-colors ${
            isButtonDisabled
              ? 'bg-gray-500 cursor-not-allowed'
              : 'bg-white hover:bg-gray-200 cursor-pointer'
          }`}
        >
          {isSubmitting ? 'Creating...' : 'Create Post'}
        </button>
      </form>
    </div>
  )
}