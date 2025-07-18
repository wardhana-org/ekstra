import CreateCreatorPageForm from "~/components/CreateCreatorPageForm";
import RotatingBackground from "~/components/RotatingBackground";
import { getBackground } from "~/server/queries";

export default async function creatorSetup() {
  const backgroundImages = await getBackground();

  return (
    <>
      <div className="flex h-screen">
        {/* Left Column */}
        <div className="relative w-full md:w-2/3">
          <div className="absolute inset-0 bg-[radial-gradient(circle_500px_at_30%_40%,#3e3e3e,transparent)] opacity-80"></div>
          <div className="relative z-10 h-full flex flex-col justify-center items-center p-8">
            <div className="max-w-md w-full rounded-xl p-8">
              <h1 className="text-4xl font-bold mb-6">Setup your creator page</h1>
              <p className="mb-6">You can edit these information later</p>
              <CreateCreatorPageForm />
            </div>
          </div>
        </div>

        {/* Right Column */}
        <div className="md:w-1/3 bg-white hidden md:block">
          <RotatingBackground images={backgroundImages} interval={4500} />
        </div>
      </div>
    </>
  );
}
