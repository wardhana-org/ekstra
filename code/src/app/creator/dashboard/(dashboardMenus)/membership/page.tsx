import CreatorMemberships from "~/components/CreatorMemberships";
import { retrieveMyMembershipTiers } from "~/server/actions";

export default async function CreatorDashboardMembership() {
  const membershipPlans = await retrieveMyMembershipTiers()
  
  return (
    <main className="p-6">
      <div className="w-1/2 mx-auto">
        <CreatorMemberships plans={membershipPlans} />
      </div>
    </main>  
  )
}