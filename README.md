# **CPM: Command-Line Tool for Smart Contract Development on EpicChain**  

**CPM** (Contract Project Manager) is a sophisticated and developer-friendly command-line utility designed for the **EpicChain** blockchain ecosystem. It enables seamless interaction with smart contracts, offering tools for development, testing, and SDK generation. Whether you're working with on-chain applications or off-chain integrations, CPM provides a robust framework to streamline your workflow, enhance productivity, and ensure a more efficient development process.  

---

## **What is CPM?**  

CPM is a multipurpose tool that simplifies complex smart contract operations for developers, reducing barriers to entry and enabling rapid prototyping and deployment. The tool integrates powerful features for interacting with smart contracts, including downloading and testing contracts from live networks and generating SDKs for a variety of programming languages.  

This utility empowers blockchain developers to:  
- Create **realistic test environments** by replicating the state of live networks within local instances.  
- Automatically generate Software Development Kits (SDKs) for diverse programming languages, helping developers integrate smart contracts with external systems easily.  
- Maintain clean, efficient, and manageable development workflows.  

---

## **Features of CPM**  

### **1. Realistic Test Environment**  
CPM allows developers to download smart contracts and their state from networks such as **MainNet** to the **EpicChain-Express** local environment. This replicates real-world conditions, ensuring accurate testing and debugging of contracts.  

### **2. Automated SDK Generation**  
CPM simplifies the process of creating SDKs by leveraging contract manifests. Developers can generate SDKs for both **on-chain** and **off-chain** interactions. Supported programming languages include:  
- **C#**: Ideal for enterprise-level blockchain integrations.  
- **Golang**: Perfect for performance-critical blockchain applications.  
- **Java**: Widely used for building decentralized applications (dApps).  
- **Python**: Great for scripting and analytics-based workflows.  
- **TypeScript**: Enables robust web-based blockchain interfaces.  

### **3. Developer-Friendly CLI Interface**  
The CLI interface is straightforward yet powerful. It provides extensive options for managing smart contracts, testing configurations, and automating SDK generation—all with simple commands.  

### **4. Cross-Platform Compatibility**  
CPM supports macOS, Windows, and Linux environments, ensuring accessibility across diverse developer ecosystems.  

### **5. Organized Output for SDKs**  
All generated SDKs are stored in a structured directory format under `/cpm_out/`. This ensures easy navigation and integration into your projects.  

---

## **Installation Instructions**  

Installing CPM is straightforward. Follow the steps below based on your operating system:  

### **Option 1: Download the Binary**  
1. Visit the [Releases Page](#).  
2. Download the latest binary for your operating system.  
3. Place the binary in a directory that is included in your system’s PATH variable.  

### **Option 2: Install via Package Managers**  

#### **macOS (Homebrew)**  
Install CPM using Homebrew, the popular macOS package manager:  
```bash
brew install epicchainlabs/tap/cpm
```  

#### **Windows (Chocolatey)**  
Install CPM using Chocolatey, a package manager for Windows:  
```bash
choco install cpm
```  

---

## **Getting Started**  

Once installed, CPM is ready to use. Begin by exploring the CLI options using the `-h` (help) flag:  
```bash
cpm -h
```  

The core configuration of your project is managed through a `cpm.yaml` file. This file defines the smart contracts and network configurations you intend to work with. Learn more about configuring your `cpm.yaml` file [here](docs/config.md).  

---

## **Usage Examples**  

### **1. Download All Contracts Listed in `cpm.yaml`**  
Download all smart contracts specified in your configuration file into your local **EpicChain-Express** environment.  
```bash
cpm --log-level DEBUG run
```  

### **2. Download a Single Contract or Manifest**  
Retrieve a specific smart contract or its manifest from a target network:  
```bash
# Download a specific contract
cpm download contract -c 0x4380f2c1de98bb267d3ea821897ec571a04fe3e0 -n mainnet

# Download only the manifest of the contract
cpm download manifest -c 0x4380f2c1de98bb267d3ea821897ec571a04fe3e0 -N https://mainnet1-seed.epic-chain.org:10111
```  

### **3. Generate SDK from Local Manifest**  
Generate on-chain or off-chain SDKs using the manifest of a smart contract:  
```bash
# Generate Python SDK for off-chain usage
cpm generate python -m samplecontract.manifest.json -t offchain

# Generate Golang SDK for on-chain usage
cpm generate go -m samplecontract.manifest.json -t onchain
```  
All SDKs will be saved in the `/cpm_out/` directory under folders specific to the SDK type and programming language. Examples:  
- **Off-chain Python SDK**: `/cpm_out/offchain/python/<contract>`  
- **On-chain Golang SDK**: `/cpm_out/onchain/golang/<contract>`  

---

## **Configuration File (`cpm.yaml`)**  

The `cpm.yaml` file is your project’s blueprint, defining which contracts to download, their target networks, and the associated settings. Learn how to create and manage this file in the [documentation](docs/config.md).  

---

## **Advanced Usage**  

For more advanced options, including setting up custom networks, integrating with existing toolchains, and optimizing workflows, refer to the [Comprehensive User Guide](docs/index.md).  

---

## **Contributions**  

We welcome contributions from the community to enhance CPM further. If you'd like to contribute, please review our [Contributing Guidelines](CONTRIBUTING.md) and join us in building the future of blockchain development.  

---

## **License**  

CPM is licensed under the [MIT License](LICENSE), ensuring its openness and accessibility for developers worldwide.  

---

## **Support and Contact**  

For support, feedback, or feature requests, please reach out via:  
- **Email**: support@epic-chain.org  
- **Community Forum**: [EpicChain Developer Community](#)  

Happy coding, and welcome to the EpicChain ecosystem! 🚀  