import CreatorDashboardHeader from "~/components/CreatorDashboardHeader";
import { retrieveCreatorPageData } from "~/server/actions";

export default async function creatorLayout({ children, }: { children: React.ReactNode; }) {

  const creatorData = await retrieveCreatorPageData();

  return(
    <>
      <div>
        <CreatorDashboardHeader creatorData={creatorData} />
        {children}
      </div>
    </>
  )
}