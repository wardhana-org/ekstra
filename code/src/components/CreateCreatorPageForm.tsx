'use client';

import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import { submitCreatorPage, checkHandleUnique } from '~/server/actions';
import { redirect } from 'next/navigation';

const formSchema = z.object({
  name: z.string().min(1, "Name is required"),
  handle: z.string()
    .min(3, "Must be at least 3 characters")
    .max(20, "Cannot exceed 20 characters"),
  description: z.string().min(10, "Description must be at least 10 characters")
});

type FormValues = z.infer<typeof formSchema>;

export default function CreateCreatorPageForm() {
  const {
    register,
    handleSubmit,
    formState: { errors, isValid },
    setError,
  } = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    mode: 'onChange'
  });

  const onSubmit = async (data: FormValues) => {
    const available = await checkHandleUnique(data.handle);
    if (!available) {
      setError("handle", {
        type: "manual",
        message: "This handle is already taken",
      });
      return;
    }

    await submitCreatorPage(data);
    redirect('creator/dashboard/home')
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      {/* Name */}
      <div>
        <label className="block pb-2">Your Name</label>
        <input
          {...register("name")}
          className="w-full px-4 py-2 rounded-lg bg-stone-800 border-2 border-gray-400 focus:outline-none focus:ring-1 focus:ring-white"
          placeholder="Enter your name"
        />
        {errors.name && (
          <p className="text-red-500 text-sm pt-1">{errors.name.message}</p>
        )}
      </div>

      {/* Handle */}
      <div>
        <label className="block pb-2">Unique Handle</label>
        <div className="relative">
          <div className="absolute left-0 top-0 bottom-0 flex items-center pl-4 pointer-events-none text-gray-400">
            domain/creatorpage/
          </div>
          <input
            {...register("handle")}
            className="w-full pl-44 px-4 py-2 rounded-lg bg-stone-800 border-2 border-gray-400 focus:outline-none focus:ring-1 focus:ring-white"
            placeholder="your-handle"
          />
        </div>
        {errors.handle && (
          <p className="text-red-500 text-sm pt-1">{errors.handle.message}</p>
        )}
      </div>

      {/* Description */}
      <div>
        <label className="block pb-2">Description</label>
        <textarea
          {...register("description")}
          className="w-full px-4 py-2 rounded-lg bg-stone-800 border-2 border-gray-400 focus:outline-none focus:ring-1 focus:ring-white"
          placeholder="Tell us about yourself"
          rows={4}
        />
        {errors.description && (
          <p className="text-red-500 text-sm pt-1">{errors.description.message}</p>
        )}
      </div>

      {/* Submit */}
      <button
        type="submit"
        disabled={!isValid}
        className={`w-full font-bold py-2 px-4 rounded-lg transition duration-200 ${
          isValid
            ? 'bg-white hover:bg-gray-200 text-black cursor-pointer'
            : 'bg-zinc-400 text-gray-800 cursor-not-allowed'
        }`}
      >
        Submit
      </button>
    </form>
  );
}
