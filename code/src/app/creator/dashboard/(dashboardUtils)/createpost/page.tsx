import CreatePostForm from "~/components/CreatePostForm";
import { retrieveMembershipNames } from "~/server/actions";


export default async function CreatePost() {
  const membership = await retrieveMembershipNames()

  return (
    <main className="m-10">    
      <div className="w-2/3 mx-auto">
        <CreatePostForm memberships={membership || []} />
      </div>
    </main>
  )
}