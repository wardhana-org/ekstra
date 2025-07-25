interface ContentRendererProps {
  content: {
    key: string;
    type: string;
  }[];
}
export default function ContentRenderer({ content }: ContentRendererProps) {
  return(
    <div></div>
  )
}