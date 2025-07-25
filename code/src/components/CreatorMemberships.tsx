'use client'

import { useRouter } from 'next/navigation';

interface MembershipPlans {
  id: string,
  title: string,
  price: string,
  description: string | null,
}

export default function CreatorMemberships( {plans} : {plans : MembershipPlans[] | null} ) {
  const router = useRouter();

  return (
    <div className="w-full p-4">
      <h2 className="text-2xl font-bold mb-6">Membership Plans</h2>
      
      <div className="space-y-4">
        {plans?.map((plan) => (
          <div 
            key={plan.id}
            className="bg-neutral-800 rounded-lg shadow-sm hover:shadow-md transition-shadow duration-200"
          >
            <div className="p-4">
              <div className="mb-2">
                <h3 className="text-lg font-semibold">{plan.title}</h3>
              </div>
              <p className="">${plan.price}<span className="text-sm">/month</span></p>
              
              {plan.description && (
                <p className="mt-2 text-sm">{plan.description}</p>
              )}
              
              <div className="mt-4 flex space-x-2">
                <button className="px-3 py-1 text-sm bg-white text-black rounded-md hover:bg-gray-200 cursor-pointer flex items-center justify-center transition-colors">
                  Edit
                </button>
              </div>
            </div>
          </div>
        ))}
      </div>
      
      <button 
        className="mt-6 w-full py-2 px-4 bg-white text-black rounded-lg hover:bg-gray-200 cursor-pointer flex items-center justify-center transition-colors"
        onClick={() => router.push('/creator/dashboard/createmembership')}        
      >
        <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 mr-2" viewBox="0 0 20 20" fill="currentColor">
          <path fillRule="evenodd" d="M10 3a1 1 0 011 1v5h5a1 1 0 110 2h-5v5a1 1 0 11-2 0v-5H4a1 1 0 110-2h5V4a1 1 0 011-1z" clipRule="evenodd" />
        </svg>
        Add New Tier
      </button>
    </div>
  );
}