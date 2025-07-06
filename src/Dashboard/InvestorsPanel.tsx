import SidePanel from './SidePanel/SidePanel';
import InvestorList from './SidePanel/InvestorList/InvestorList';

interface InvestorsPanelProps {
  visible: boolean;
}

export default function InvestorsPanel({ visible }: InvestorsPanelProps) {
  return (
    <SidePanel side="left" visible={visible}>
      <InvestorList />
    </SidePanel>
  );
}
