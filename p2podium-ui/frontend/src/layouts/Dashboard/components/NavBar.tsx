import P2Podium from "../../../assets/P2Podium.svg";

export default function Header() {
  return (
    <header className="fixed top-5 left-10">
      <img src={P2Podium} alt="P2Podium Logo" className="h-[120px]" />
    </header>
  );
}
