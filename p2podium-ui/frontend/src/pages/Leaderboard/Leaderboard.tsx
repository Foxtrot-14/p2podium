import { useRouter } from "../../routes/hooks/use-router";

export default function Leaderboard() {
  const navigate = useRouter();

  const OnClick = () => {
    navigate.push("/")
  }

  return (
    <article className="h-fit w-fit">
      <h1 className="text-[#F90627]">This is leaderboard</h1>
      <button onClick={OnClick}>Go Back</button>
    </article>
  );
}
