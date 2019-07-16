package i3tech.chaincode;


import static org.assertj.core.api.AssertionsForClassTypes.assertThat;
import static org.mockito.BDDMockito.given;

import org.hyperledger.fabric.shim.Chaincode;
import org.hyperledger.fabric.shim.ChaincodeStub;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.Mock;
import org.mockito.runners.MockitoJUnitRunner;

@RunWith(MockitoJUnitRunner.class)
public class SmartContractTest {

    @Mock
    private ChaincodeStub chaincodeStub;
    private SmartContract chaincode = new SmartContract();

    @Test
    public void shouldReturnErrorForIncorrectFunctionName() {
        //given
        String functionName = "wrong_name";
        given(chaincodeStub.getFunction()).willReturn(functionName);

        //when
        Chaincode.Response result = chaincode.invoke(chaincodeStub);

        //then
        assertThat(result.getStatusCode()).isEqualTo(500);
        assertThat(result.getMessage()).isEqualTo("wrong_name function is currently not supported");
    }
}