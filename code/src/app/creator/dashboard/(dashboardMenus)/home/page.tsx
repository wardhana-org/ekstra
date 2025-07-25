import CreatorHome from "~/components/CreatorHome";
import Posts from "~/components/dummyData/DummyData"

export default function CreatorDashboardHome() {
  return (
    <main className="p-6">
      <div className="w-1/2 mx-auto">
        <CreatorHome posts={Posts} />
      </div>
    </main>  
  )
}