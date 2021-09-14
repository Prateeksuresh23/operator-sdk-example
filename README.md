# operator-sdk-example

This repo contains example of operator which take image name from CR and create deployment. For creating operator used operator-sdk framework. For debugging and faster development using okteto.

# Install operator sdk
RELEASE_VERSION=v1.1.0 <br />
curl -LO https://github.com/operator-framework/operator-sdk/releases/download/${RELEASE_VERSION}/operator-sdk-${RELEASE_VERSION}-x86_64-linux-gnu <br />

chmod +x operator-sdk-${RELEASE_VERSION}-x86_64-linux-gnu && sudo mkdir -p /usr/local/bin/ && sudo cp operator-sdk-${RELEASE_VERSION}-x86_64-linux-gnu /usr/local/bin/operator-sdk && rm operator-sdk-${RELEASE_VERSION}-x86_64-linux-gnu

operator-sdk init --domain example.com <br />
tree -d . <br />
operator-sdk create api --resource --controller --group hello --kind TestOp --version=v1beta1 <br />
make manifests <br />

export IMG = prateeksuresh23/operator-example:0.1.0 <br />
make docker-build docker-push IMG=$IMG <br />
make deploy IMG=$IMG <br />

# okteto
okteto sync our local code with application running on containers. Using this we need not create new image after every changes <br />
curl https://get.okteto.com -sSfL | sh <br />
okteto init -n=operator-system <br />
okteto up -n=operator-system <br />