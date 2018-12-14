require 'spec_helper_acceptance'

describe '::site_module::containerd' do
  context 'containerd' do
    let :pp do
      "
class{'site_module::containerd':
}
"
    end

    before(:all) do
      hosts.each do |host|
        # reset firewall
        on host, "iptables -F INPUT"

        # ensure no swap space is mounted
        on host, "swapoff -a"
      end
    end

    it 'should converge on the first puppet run' do
      hosts.each do |host|
        apply_manifest_on(host, pp, :catch_failures => true)
        expect(
          apply_manifest_on(host, pp, :catch_failures => true).exit_code
        ).to be_zero
      end
    end
  end
end
