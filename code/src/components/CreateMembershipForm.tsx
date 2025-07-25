'use client';

import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import { submitMembershipTierInfo } from '~/server/actions';
import { useRouter } from 'next/navigation';
import { useState } from 'react';

const formSchema = z.object({
  title: z.string().min(1, "Title is required"),
  price: z.number().min(1, "Price must be at least $1").max(1000, "Price cannot exceed $1000"),
  description: z.string().min(10, "Description must be at least 10 characters")
});

type FormValues = z.infer<typeof formSchema>;

export default function CreateMembershipForm() {
  const [isProcessing, setIsProcessing] = useState(false);
  
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    mode: 'onChange'
  });

  const router = useRouter();
  const onSubmit = async (data: FormValues) => {
    setIsProcessing(true);
    try {
      await submitMembershipTierInfo(data)
      router.push('/creator/dashboard/membership')
      alert('Membership created successfully!');
    } catch (error) {
      console.error('Error creating membership:', error);
      alert('Failed to create membership');
      setIsProcessing(false);
    }
  };

  const isButtonDisabled = isSubmitting || isProcessing;

  const labelStyle = 'block pb-2'
  const inputStyle = 'w-full px-4 py-2 rounded-lg bg-stone-800 border-2 border-gray-400 focus:outline-none focus:ring-1 focus:ring-white'
 
  return (
    <div className="w-full p-4 bg-neutral-800 rounded-lg shadow-sm">
      <h2 className="text-2xl font-bold mb-6">Membership Details</h2>
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
        {/* Membership Name */}
        <div>
          <label className={labelStyle}>
            Membership Name
          </label>
          <input
            {...register("title")}
            className={inputStyle}
            placeholder="Bronze Tier"
          />
          {errors.title && (
            <p className="mt-1 text-sm text-red-500">{errors.title.message}</p>
          )}
        </div>

        {/* Monthly Cost */}
        <div>
          <label className={labelStyle}>
            Monthly Cost ($)
          </label>
          <input
            type="number"
            {...register("price", { valueAsNumber: true })}
            className={inputStyle}
            placeholder="15.00"
            step="0.01"
          />
          {errors.price && (
            <p className="mt-1 text-sm text-red-500">{errors.price.message}</p>
          )}
        </div>

        {/* Membership Description */}
        <div>
          <label className={labelStyle}>
            Membership Description
          </label>
          <textarea
            {...register("description")}
            className="w-full px-4 py-2 rounded-lg bg-stone-800 border-2 border-gray-400 focus:outline-none focus:ring-1 focus:ring-white"
            placeholder={`- Access to behind the scenes content\n- Exclusive episodes\n- Early access to new content\n- Members-only community`}
            rows={5}
          />
          {errors.description && (
            <p className="mt-1 text-sm text-red-500">{errors.description.message}</p>
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
          {isSubmitting ? 'Creating...' : 'Create Membership'}
        </button>
      </form>
    </div>
  );
}