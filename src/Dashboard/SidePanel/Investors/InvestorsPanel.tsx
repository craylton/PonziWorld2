import SidePanel from '../SidePanel';
import InvestorList from './InvestorList';

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
